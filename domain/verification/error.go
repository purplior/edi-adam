package verification

import "errors"

var (
	ErrInvalidCode     = errors.New("invalid code")
	ErrAlreadyVerified = errors.New("already verified")
)
