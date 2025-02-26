package service

import "errors"

var (
	ErrTooManyRetries = errors.New("too many retries")
	ErrInvalidUUID    = errors.New("invalid uuid")
)
