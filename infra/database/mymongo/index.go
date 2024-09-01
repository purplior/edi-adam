package mymongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitIndexes(
	ctx context.Context,
	client *Client,
) error {
	if _, err := client.
		MyCollection(Collection_User).
		Indexes().
		CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: "join_method", Value: 1},
					{Key: "account_id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		}); err != nil {
		return err
	}

	return nil
}
