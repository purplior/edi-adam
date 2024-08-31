package mymongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type (
	Transaction struct {
		session mongo.Session
	}
)

func (s *Transaction) Run(
	task func(ctx context.Context) (interface{}, error),
	opts ...*options.TransactionOptions,
) (result interface{}, err error) {
	baseOptions := []*options.TransactionOptions{
		options.
			Transaction().
			SetWriteConcern(&writeconcern.WriteConcern{
				W: "majority",
			}),
	}
	baseOptions = append(baseOptions, opts...)

	return s.session.WithTransaction(
		context.TODO(),
		func(ctx mongo.SessionContext) (interface{}, error) {
			return task(ctx)
		},
		baseOptions...,
	)
}

func NewMyMongoTransaction(
	session mongo.Session,
) *Transaction {
	return &Transaction{
		session: session,
	}
}
