package entity

import (
	"github.com/google/uuid"
	"slices"
	"strings"
	"time"
)

type Wallet struct {
	UUID      uuid.UUID
	Amount    int64
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWallet() *Wallet {
	return &Wallet{
		Amount:    0,
		Version:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (w *Wallet) DoTransaction(t *Transaction) (*Wallet, error) {
	if err := t.isValid(); err != nil {
		return nil, err
	}
	if w.UUID != t.WalletUUID {
		return nil, ErrInvalidOperationUUID
	}

	copyWallet := &Wallet{
		UUID:      w.UUID,
		Amount:    w.Amount,
		Version:   w.Version,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}

	switch t.Operation {
	case Withdraw:
		copyWallet.withdraw(t)
	case Deposit:
		copyWallet.deposit(t)
	default:
		return nil, ErrInvalidOperationType
	}

	return copyWallet, copyWallet.validate()
}

func (w *Wallet) withdraw(t *Transaction) {
	w.Amount -= t.Amount
	w.UpdatedAt = time.Now()
}

func (w *Wallet) deposit(t *Transaction) {
	w.Amount += t.Amount
	w.UpdatedAt = time.Now()
}

func (w *Wallet) validate() error {
	if w.Amount < 0 {
		return ErrNotEnoughFunds
	}
	return nil
}

type Transaction struct {
	ID             int64
	WalletUUID     uuid.UUID
	IdempotencyKey uuid.UUID
	Operation      OperationType
	Amount         int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OperationType string

var (
	Withdraw OperationType = "withdraw"
	Deposit  OperationType = "deposit"
)

func NewOperation(walletUUID uuid.UUID, operationType string, amount int64) (Transaction, error) {
	operationType = strings.ToLower(operationType)

	o := Transaction{
		WalletUUID:     walletUUID,
		IdempotencyKey: uuid.New(),
		Operation:      OperationType(operationType),
		Amount:         amount,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return o, o.isValid()
}

func (t *Transaction) isValid() error {
	if t.Amount <= 0 {
		return ErrAmountIsOrBelowZero
	}
	if t.WalletUUID == uuid.Nil {
		return ErrWalletUUIDIsEmpty
	}
	if !slices.Contains([]OperationType{Withdraw, Deposit}, t.Operation) {
		return ErrInvalidOperationType
	}
	return nil
}
