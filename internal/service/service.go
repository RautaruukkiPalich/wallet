package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
	"sync"
	"wallet/internal/entity"
	transactionRepository "wallet/internal/repository/transaction"
	walletRepository "wallet/internal/repository/wallet"
)

//go:generate mockery --name walletRepo --structname=WalletRepo
type walletRepo interface {
	Insert(context.Context, pgx.Tx, *entity.Wallet) error
	Update(context.Context, pgx.Tx, *entity.Wallet) error

	GetByUUID(context.Context, pgx.Tx, uuid.UUID) (*entity.Wallet, error)
}

//go:generate mockery --name walletCache --structname=WalletCache
type walletCache interface {
	SetBalance(context.Context, uuid.UUID, int64) error
	GetBalance(context.Context, uuid.UUID) (int64, error)
}

//go:generate mockery --name transactionRepo --structname=TransactionRepo
type transactionRepo interface {
	Insert(context.Context, pgx.Tx, *entity.Transaction) error
	Exists(context.Context, pgx.Tx, *entity.Transaction) (bool, error)
}

//go:generate mockery --name transactionBroker --structname=TransactionBroker
type transactionBroker interface {
	Publish(context.Context, *entity.Transaction) error
	Consume(context.Context) (*entity.Transaction, error)
}

//go:generate mockery --name store --structname=Store
type store interface {
	WithTransact(context.Context, func(pgx.Tx) error) error
	//WithSerializableTransact(context.Context, func(pgx.Tx) error) error
}

type Service struct {
	walletRepo        walletRepo
	transactionRepo   transactionRepo
	transactionBroker transactionBroker
	walletCache       walletCache
	store             store
	workersCount      int8
	mu                *sync.RWMutex
	walletMutex       sync.Map
}

func New(
	ctx context.Context,
	walletRepo walletRepo,
	transactionRepo transactionRepo,
	transactionBroker transactionBroker,
	walletCache walletCache,
	store store,
	workersCount int8,

) *Service {
	s := &Service{
		walletRepo:        walletRepo,
		transactionRepo:   transactionRepo,
		transactionBroker: transactionBroker,
		walletCache:       walletCache,
		store:             store,
		workersCount:      workersCount,
		mu:                &sync.RWMutex{},
	}

	for range s.workersCount {
		go s.consumeTransactions(ctx)
	}

	return s
}

/*
NEW WALLET
*/

func (s *Service) NewWallet(ctx context.Context) (*entity.Wallet, error) {
	wallet := entity.NewWallet()

	err := s.store.WithTransact(ctx, func(t pgx.Tx) error {
		return s.walletRepo.Insert(ctx, t, wallet)
	})
	if err != nil {
		return nil, err
	}

	return wallet, err
}

/*
GET BALANCE
*/

func (s *Service) GetBalance(ctx context.Context, uid uuid.UUID) (int64, error) {

	var balance int64
	var err error

	if uuid.Nil == uid {
		return 0, ErrInvalidUUID
	}

	if err = s.withRLock(func() error {
		balance, err = s.walletCache.GetBalance(ctx, uid)
		return err
	}); err == nil {
		return balance, nil
	}

	err = s.withLock(func() error {
		return s.updateCache(ctx, uid)
	})
	if err != nil {
		return 0, err
	}

	balance, err = s.walletCache.GetBalance(ctx, uid)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (s *Service) updateCache(ctx context.Context, uid uuid.UUID) error {
	return s.store.WithTransact(ctx, func(tx pgx.Tx) error {
		wallet, err := s.walletRepo.GetByUUID(ctx, tx, uid)
		if err != nil {
			return err
		}

		err = s.walletCache.SetBalance(ctx, uid, wallet.Amount)
		if err != nil {
			log.Println("error update cache: ", err)
		}

		return nil
	})
}

func (s *Service) withLock(fn func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fn()
}

func (s *Service) withRLock(fn func() error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fn()
}

/*
TRANSACTION WITH BROKER
*/

func (s *Service) NewTransaction(ctx context.Context, t *entity.Transaction) error {
	err := s.store.WithTransact(ctx, func(tx pgx.Tx) error {
		wallet, err := s.walletRepo.GetByUUID(ctx, tx, t.WalletUUID)
		if err != nil {
			return err
		}
		_, err = wallet.DoTransaction(t)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, entity.ErrWalletUUIDIsEmpty) {
			return err
		}
		if errors.Is(err, entity.ErrNotEnoughFunds) {
			return err
		}
		if errors.Is(err, walletRepository.ErrWalletNotFound) {
			return err
		}

		log.Println("error while get wallet wallet:", err)
		return err
	}

	return s.transactionBroker.Publish(ctx, t)
}

func (s *Service) consumeTransactions(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t, err := s.transactionBroker.Consume(ctx)
			if err != nil {
				log.Println("error consuming transaction: ", err)
				continue
			}

			mu := s.getWalletMutex(t.WalletUUID)

			mu.Lock()

			err = s.store.WithTransact(ctx, func(tx pgx.Tx) error {
				if exists, err := s.transactionRepo.Exists(ctx, tx, t); err != nil || exists {
					return nil
				}
				wallet, err := s.walletRepo.GetByUUID(ctx, tx, t.WalletUUID)
				if err != nil {
					return err
				}
				newWallet, err := wallet.DoTransaction(t)
				if err != nil {
					return err
				}
				err = s.walletRepo.Update(ctx, tx, newWallet)
				if err != nil {
					return err
				}

				t.StatusSuccess()
				err = s.transactionRepo.Insert(ctx, tx, t)
				if err != nil {
					return err
				}

				err = s.walletCache.SetBalance(ctx, t.WalletUUID, newWallet.Amount)
				if err != nil {
					log.Println("error set cache: ", err)
				}

				return nil
			})

			mu.Unlock()

			if err != nil {
				log.Printf("failed to process transaction %v: %v\n", t.IdempotencyKey, err)

				if errors.Is(err, transactionRepository.ErrDuplicateTransaction) ||
					errors.Is(err, walletRepository.ErrWalletNotFound) {
					continue
				}

				s.handleTransactionError(ctx, t, err)
			}
		}
	}
}

func (s *Service) handleTransactionError(ctx context.Context, t *entity.Transaction, err error) {
	if errors.Is(err, transactionRepository.ErrDuplicateTransaction) {
		log.Printf("Duplicate transaction: %v, skipping\n", t.IdempotencyKey)
		return
	}

	if errors.Is(err, entity.ErrNotEnoughFunds) {
		s.markTransactionAsFailed(ctx, t)
		return
	}

	t.StatusNew()
	if err := s.transactionBroker.Publish(ctx, t); err != nil {
		log.Printf("Failed to requeue transaction %v: %v\n", t.IdempotencyKey, err)
	}
}

func (s *Service) markTransactionAsFailed(ctx context.Context, t *entity.Transaction) {
	err := s.store.WithTransact(ctx, func(tx pgx.Tx) error {
		t.StatusFailure()
		return s.transactionRepo.Insert(ctx, tx, t)
	})
	if err != nil {
		log.Printf("Failed to mark transaction as failed: %v", err)
	}
}

func (s *Service) getWalletMutex(uid uuid.UUID) *sync.Mutex {
	mu, _ := s.walletMutex.LoadOrStore(uid.String(), &sync.Mutex{})
	return mu.(*sync.Mutex)
}
