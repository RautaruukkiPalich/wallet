package transaction

import "errors"

var (
	ErrDuplicateTransaction = errors.New("duplicate transaction")
)
