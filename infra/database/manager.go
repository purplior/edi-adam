package database

import (
	"context"

	"github.com/podossaem/root/infra/database/mymongo"
)

func Init(
	client *mymongo.Client,
) error {
	ctx := context.Background()

	if err := client.Connect(ctx); err != nil {
		return err
	}

	return nil
}

func Dispose(client *mymongo.Client) error {
	return client.Disconnect(context.Background())
}
