package exception

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoDefaultDB = errors.New("no default db")
	ErrNoDocuments = mongo.ErrNoDocuments
)
