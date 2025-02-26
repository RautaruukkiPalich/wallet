package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sync"
	"wallet/internal/entity"
	walletRepository "wallet/internal/repository/wallet"
)

type walletRepo interface {
	Insert(context.Context, pgx.Tx, *entity.Wallet) error
	Update(context.Context, pgx.Tx, *entity.Wallet) error

	GetByUUID(context.Context, pgx.Tx, uuid.UUID) (*entity.Wallet, error)
}

type walletCache interface {
	SetBalance(context.Context, uuid.UUID, int64) error
	GetBalance(context.Context, uuid.UUID) (int64, error)
}

type transactionRepo interface {
	Insert(context.Context, pgx.Tx, *entity.Transaction) error
}

type store interface {
	WithTransact(context.Context, func(pgx.Tx) error) error
	WithSerializableTransact(context.Context, func(pgx.Tx) error) error
}

type Service struct {
	walletRepo      walletRepo
	transactionRepo transactionRepo
	walletCache     walletCache
	store           store
	retry           int8
	mu              *sync.RWMutex
}

func New(
	walletRepo walletRepo,
	transactionRepo transactionRepo,
	walletCache walletCache,
	store store,
	retry int8,
) *Service {
	return &Service{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		walletCache:     walletCache,
		store:           store,
		retry:           retry,
		mu:              &sync.RWMutex{},
	}
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

	wallet := new(entity.Wallet)

	err = s.withLock(func() error {
		return s.store.WithTransact(ctx, func(tx pgx.Tx) error {

			wallet, err = s.walletRepo.GetByUUID(ctx, tx, uid)
			if err != nil {
				return err
			}
			return s.walletCache.SetBalance(ctx, uid, wallet.Amount)

		})
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
NEW TRANSACTION
*/

func (s *Service) NewTransaction(ctx context.Context, t *entity.Transaction) error {
	var balance int64

	for range s.retry {
		err := s.store.WithSerializableTransact(ctx, func(tx pgx.Tx) error {
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
			err = s.transactionRepo.Insert(ctx, tx, t)
			if err != nil {
				return err
			}
			balance = newWallet.Amount
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
			continue
		}

		if err := s.walletCache.SetBalance(ctx, t.WalletUUID, balance); err != nil {
			fmt.Println("error set cache: ", err)
		}

		return nil
	}

	return ErrTooManyRetries
}
