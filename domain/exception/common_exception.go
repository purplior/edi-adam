package exception

import "errors"

var (
	ErrBadRequest    = errors.New("bad request")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrNotAcceptable = errors.New("not acceptable")
)
