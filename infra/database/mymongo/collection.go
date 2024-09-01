package mymongo

import "go.mongodb.org/mongo-driver/mongo"

const (
	Collection_EmailVerification = "email_verifications"
	Collection_User              = "users"
)

type (
	MyMongoCollection struct {
		*mongo.Collection
	}
)

func NewCollection(
	collection *mongo.Collection,
) *MyMongoCollection {
	return &MyMongoCollection{
		Collection: collection,
	}
}
