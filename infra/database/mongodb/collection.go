package mongodb

import "go.mongodb.org/mongo-driver/mongo"

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
