package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
	"wallet/internal/entity"
)

type transactor interface {
	WithTransact(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type walletRepo interface {
	Create(ctx context.Context, tx pgx.Tx, wallet *entity.Wallet) error
	Update(ctx context.Context, tx pgx.Tx, wallet *entity.Wallet) error
	GetByUUID(ctx context.Context, tx pgx.Tx, walletUUID uuid.UUID) (*entity.Wallet, error)

	GetBalanceByUUID(ctx context.Context, walletUUID uuid.UUID) (int64, error)
	SetBalanceByUUID(ctx context.Context, wallet *entity.Wallet) error
}

type transactionRepo interface {
	Create(ctx context.Context, tx pgx.Tx, transaction *entity.Transaction) error
}

type WalletService struct {
	transactor      transactor
	walletRepo      walletRepo
	transactionRepo transactionRepo
}

func NewWalletService() *WalletService {
	return &WalletService{}
}

func (s *WalletService) GetBalance(ctx context.Context, uid uuid.UUID) (int64, error) {
	var balance int64

	if uid == uuid.Nil {
		return 0, entity.ErrWalletUUIDIsEmpty
	}

	cachedBalance, err := s.walletRepo.GetBalanceByUUID(ctx, uid)
	if err == nil {
		return cachedBalance, nil
	}

	if err := s.transactor.WithTransact(ctx, func(tx pgx.Tx) error {
		wallet, err := s.walletRepo.GetByUUID(ctx, tx, uid)
		if err != nil {
			return err
		}

		balance = wallet.Amount
		s.setCacheBalanceByWallet(ctx, wallet)
		return nil
	}); err != nil {
		return 0, err
	}

	return balance, nil
}

const maxRetries = 5

func (s *WalletService) NewTransaction(ctx context.Context, t *entity.Transaction) error {
	var err error
	var wallet *entity.Wallet

	for range maxRetries {
		err = s.transactor.WithTransact(ctx, func(tx pgx.Tx) error {
			wallet, err = s.walletRepo.GetByUUID(ctx, tx, t.WalletUUID)
			if err != nil {
				return err
			}

			if err = wallet.DoTransaction(t); err != nil {
				return err
			}

			if err = s.transactionRepo.Create(ctx, tx, t); err != nil {
				return err
			}

			if err = s.walletRepo.Update(ctx, tx, wallet); err != nil {
				return err
			}
			return nil
		})

		if err == nil {
			s.setCacheBalanceByWallet(ctx, wallet)
			return nil
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *WalletService) setCacheBalanceByWallet(ctx context.Context, wallet *entity.Wallet) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if cacheErr := s.walletRepo.SetBalanceByUUID(ctx, wallet); cacheErr != nil {
		fmt.Printf("fail to set cache for wallet %s: %v\n", wallet.UUID, cacheErr)
	}
}
