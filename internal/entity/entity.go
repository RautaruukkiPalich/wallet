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
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWallet() Wallet {
	return Wallet{
		Amount:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (w *Wallet) DoTransaction(t *Transaction) error {
	if err := t.isValid(); err != nil {
		return err
	}
	if w.UUID != t.WalletUUID {
		return ErrInvalidOperationUUID
	}

	switch t.Operation {
	case Withdraw:
		w.withdraw(t)
	case Deposit:
		w.deposit(t)
	default:
		return ErrInvalidOperationType
	}

	return w.validate()
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
	WalletUUID uuid.UUID
	Operation  OperationType
	Amount     int64
	CreatedAt  time.Time
}

type OperationType string

var (
	Withdraw OperationType = "withdraw"
	Deposit  OperationType = "deposit"
)

func NewOperation(walletUUID uuid.UUID, operationType string, amount int64) (Transaction, error) {
	operationType = strings.ToLower(operationType)

	o := Transaction{
		WalletUUID: walletUUID,
		Operation:  OperationType(operationType),
		Amount:     amount,
		CreatedAt:  time.Now(),
	}

	return o, o.isValid()
}

func (t Transaction) isValid() error {
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
