package exception

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoDefaultDB   = errors.New("no default db")
	ErrNoDocuments   = mongo.ErrNoDocuments
	ErrNoRecord      = errors.New("no record")
	ErrDBProcess     = errors.New("db process")
	ErrInTransaction = errors.New("tx in tx")
	ErrNoTransaction = errors.New("no tx")
)
