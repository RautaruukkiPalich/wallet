package httpserver

import "errors"

var (
	ErrDuplicateRun = errors.New("duplicate http server running")
)
