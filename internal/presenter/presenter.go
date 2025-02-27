package presenter

import (
	"context"
	"github.com/google/uuid"
	"wallet/internal/dto"
	"wallet/internal/entity"
)

type walletService interface {
	NewTransaction(ctx context.Context, operation *entity.Transaction) error

	NewWallet(ctx context.Context) (*entity.Wallet, error)
	GetBalance(context.Context, uuid.UUID) (int64, error)
}

type Presenter struct {
	walletService walletService
}

func NewPresenter(walletService walletService) *Presenter {
	return &Presenter{
		walletService: walletService,
	}
}

func (p *Presenter) Transaction(ctx context.Context, req *dto.PostOperationRequest) error {
	walletUUID, err := uuid.Parse(req.WalletId)
	if err != nil {
		return ErrInvalidUUID
	}

	operation, err := entity.NewOperation(walletUUID, req.OperationType, req.Amount)
	if err != nil {
		return err
	}

	if err := p.walletService.NewTransaction(ctx, operation); err != nil {
		return err
	}

	return nil
}

func (p *Presenter) GetBalance(ctx context.Context, uid string) (*dto.GetBalanceResponse, error) {
	walletUUID, err := uuid.Parse(uid)
	if err != nil || walletUUID == uuid.Nil {
		return nil, ErrInvalidUUID
	}

	balance, err := p.walletService.GetBalance(ctx, walletUUID)
	if err != nil {
		return nil, err
	}

	return &dto.GetBalanceResponse{Amount: balance}, nil
}

func (p *Presenter) NewWallet(ctx context.Context) (*dto.WalletResponse, error) {
	wallet, err := p.walletService.NewWallet(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.WalletResponse{UUID: wallet.UUID.String(), Amount: wallet.Amount}, nil
}
