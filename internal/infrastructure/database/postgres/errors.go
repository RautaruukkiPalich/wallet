package postgres

import "errors"

var (
	ErrConnectToDB = errors.New("error connecting to db")
	ErrPingToDB    = errors.New("error ping to db")
)
