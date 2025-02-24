package api

import "errors"

var (
	ErrEmptyWalletUUID = errors.New("empty wallet uuid")
	ErrInvalidFormData = errors.New("invalid form data")
)
