package exception

import "errors"

var (
	ErrInvalidVerificationCode = errors.New("invalid code")
	ErrAlreadyVerified         = errors.New("already verified")
	ErrAlreadyConsumed         = errors.New("already consumed")
)
