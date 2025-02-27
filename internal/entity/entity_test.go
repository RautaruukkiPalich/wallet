package entity

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()

	if wallet.Amount != 0 {
		t.Errorf("Expected Amount to be 0, got %d", wallet.Amount)
	}
	if wallet.Version != 0 {
		t.Errorf("Expected Version to be 0, got %d", wallet.Version)
	}
	if wallet.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set, got zero value")
	}
	if wallet.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set, got zero value")
	}
}

func TestWallet_Deposit(t *testing.T) {
	wallet := NewWallet()
	wallet.UUID = uuid.New()

	transaction, err := NewOperation(
		wallet.UUID,
		"DEPOSIT",
		100,
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	updatedWallet, err := wallet.DoTransaction(transaction)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedWallet.Amount != 100 {
		t.Errorf("Expected Amount to be 100, got %d", updatedWallet.Amount)
	}
}

func TestWallet_Withdraw(t *testing.T) {
	wallet := NewWallet()
	wallet.UUID = uuid.New()
	wallet.Amount = 200

	transaction, err := NewOperation(
		wallet.UUID,
		"WITHDRAW",
		100,
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	updatedWallet, err := wallet.DoTransaction(transaction)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedWallet.Amount != 100 {
		t.Errorf("Expected Amount to be 100, got %d", updatedWallet.Amount)
	}
}

func TestWallet_Withdraw_InsufficientFunds(t *testing.T) {
	wallet := NewWallet()
	wallet.Amount = 50

	transaction := &Transaction{
		WalletUUID: wallet.UUID,
		Operation:  Withdraw,
		Amount:     100,
	}

	_, err := wallet.DoTransaction(transaction)
	if err == nil {
		t.Error("Expected error for insufficient funds, got nil")
	}
}

func TestNewOperation(t *testing.T) {
	walletUUID := uuid.New()
	transaction, err := NewOperation(walletUUID, "deposit", 100)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if transaction.WalletUUID != walletUUID {
		t.Errorf("Expected WalletUUID to be %v, got %v", walletUUID, transaction.WalletUUID)
	}
	if transaction.Operation != Deposit {
		t.Errorf("Expected Operation to be Deposit, got %v", transaction.Operation)
	}
	if transaction.Amount != 100 {
		t.Errorf("Expected Amount to be 100, got %d", transaction.Amount)
	}
	if transaction.Status != New {
		t.Errorf("Expected Status to be New, got %v", transaction.Status)
	}
}

func TestInvalidNewOperation(t *testing.T) {
	walletUUID := uuid.New()
	_, err := NewOperation(walletUUID, "sadsadsa", 100)

	if !errors.Is(err, ErrInvalidOperationType) {
		t.Error("Expected ErrInvalidOperationType")
	}
}

func TestTransaction_IsValid(t *testing.T) {
	tcs := []struct {
		name        string
		transaction *Transaction
		expectedErr error
	}{
		{
			name: "Valid transaction",
			transaction: &Transaction{
				WalletUUID: uuid.New(),
				Operation:  Deposit,
				Amount:     100,
				Status:     New,
			},
			expectedErr: nil,
		},
		{
			name: "Invalid operation type",
			transaction: &Transaction{
				WalletUUID: uuid.New(),
				Operation:  "invalid",
				Amount:     100,
				Status:     New,
			},
			expectedErr: ErrInvalidOperationType,
		},
		{
			name: "Invalid status",
			transaction: &Transaction{
				WalletUUID: uuid.New(),
				Operation:  Deposit,
				Amount:     100,
				Status:     "invalid",
			},
			expectedErr: ErrInvalidStatus,
		},
		{
			name: "Amount is zero",
			transaction: &Transaction{
				WalletUUID: uuid.New(),
				Operation:  Deposit,
				Amount:     0,
				Status:     New,
			},
			expectedErr: ErrAmountIsOrBelowZero,
		},
		{
			name: "WalletUUID is empty",
			transaction: &Transaction{
				WalletUUID: uuid.Nil,
				Operation:  Deposit,
				Amount:     100,
				Status:     New,
			},
			expectedErr: ErrWalletUUIDIsEmpty,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.transaction.isValid()
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestTransaction_StatusChange(t *testing.T) {
	transaction := &Transaction{
		Status: New,
	}

	transaction.StatusSuccess()
	if transaction.Status != Success {
		t.Errorf("Expected Status to be Success, got %v", transaction.Status)
	}

	transaction.StatusFailure()
	if transaction.Status != Failure {
		t.Errorf("Expected Status to be Failure, got %v", transaction.Status)
	}

	transaction.StatusNew()
	if transaction.Status != New {
		t.Errorf("Expected Status to be New, got %v", transaction.Status)
	}
}

func TestTransaction_Marshall(t *testing.T) {
	transaction := &Transaction{
		WalletUUID:     uuid.New(),
		IdempotencyKey: uuid.New(),
		Operation:      Deposit,
		Amount:         100,
		Status:         New,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	data, err := transaction.Marshall()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty data, got empty")
	}
}

func TestTransaction_Unmarshall(t *testing.T) {
	transaction := &Transaction{
		WalletUUID:     uuid.New(),
		IdempotencyKey: uuid.New(),
		Operation:      Deposit,
		Amount:         100,
		Status:         New,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	data, err := transaction.Marshall()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	newTransaction := &Transaction{}
	err = newTransaction.Unmarshall(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if newTransaction.WalletUUID != transaction.WalletUUID {
		t.Errorf("Expected WalletUUID to be %v, got %v", transaction.WalletUUID, newTransaction.WalletUUID)
	}
}
