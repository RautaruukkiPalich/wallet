package wallet

import "errors"

var (
	ErrWalletNotFound = errors.New("wallet not found")
	ErrNoRowsAffected = errors.New("no rows affected")
)
