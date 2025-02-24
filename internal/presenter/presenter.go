package presenter

import (
	"context"
	"github.com/google/uuid"
	"wallet/internal/dto"
	"wallet/internal/entity"
)

type walletService interface {
	NewTransaction(ctx context.Context, operation *entity.Transaction) error
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
	walletUUID, err := uuid.Parse(req.WalletUUID)
	if err != nil {
		return err
	}

	operation, err := entity.NewOperation(walletUUID, req.OperationType, req.Amount)
	if err != nil {
		return err
	}

	if err := p.walletService.NewTransaction(ctx, &operation); err != nil {
		return err
	}

	return nil
}

func (p *Presenter) GetBalance(ctx context.Context, uid string) (*dto.GetBalanceResponse, error) {
	walletUUID, err := uuid.Parse(uid)
	if err != nil {
		return nil, err
	}

	balance, err := p.walletService.GetBalance(ctx, walletUUID)
	if err != nil {
		return nil, err
	}

	return &dto.GetBalanceResponse{Amount: balance}, nil
}
