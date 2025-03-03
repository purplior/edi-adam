package exception

import (
	"errors"
)

var (
	ErrNoDefaultDB        = errors.New("no default db")
	ErrNoRecord           = errors.New("no record")
	ErrDBProcess          = errors.New("db process")
	ErrInvalidTransaction = errors.New("invalid transaction")
	ErrInTransaction      = errors.New("in transaction")
	ErrNoTransaction      = errors.New("no transaction")
	ErrNoAffected         = errors.New("no affected")
)
