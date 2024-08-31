package mymongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MyMongoDatabase struct {
		*mongo.Database
		collectionMap map[string]*MyMongoCollection
	}
)

func (d *MyMongoDatabase) MyCollection(
	name string,
	opts ...*options.CollectionOptions,
) (collection *MyMongoCollection) {
	col, is := d.collectionMap[name]
	if !is || col == nil {
		d.collectionMap[name] = NewCollection(
			d.Collection(name, opts...),
		)
		col = d.collectionMap[name]
	}

	return col
}

func NewDatabase(
	db *mongo.Database,
) *MyMongoDatabase {
	return &MyMongoDatabase{
		Database:      db,
		collectionMap: make(map[string]*MyMongoCollection),
	}
}
