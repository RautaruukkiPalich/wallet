package entity

import "errors"

var (
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrNotEnoughFunds       = errors.New("not enough funds")
	ErrWalletUUIDIsEmpty    = errors.New("wallet uuid is empty")
	ErrInvalidOperationUUID = errors.New("invalid operation uuid")
	ErrAmountIsOrBelowZero  = errors.New("amount is or below zero")
)
