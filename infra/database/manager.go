package database

import (
	"context"

	"github.com/podossaem/podoroot/infra/database/mymongo"
	"github.com/podossaem/podoroot/infra/database/myredis"
)

func Init(
	mymongoClient *mymongo.Client,
	myredisClient *myredis.Client,
) error {
	ctx := context.Background()

	if err := mymongoClient.Connect(ctx); err != nil {
		return err
	}
	if err := mymongo.InitIndexes(ctx, mymongoClient); err != nil {
		return err
	}

	// 현재 Redis 사용하는 곳이 없어서 주석
	// if err := myredisClient.Connect(ctx); err != nil {
	// 	return err
	// }

	return nil
}

func Dispose(client *mymongo.Client) error {
	return client.Disconnect(context.Background())
}
