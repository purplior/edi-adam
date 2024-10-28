package repository

import (
	"github.com/podossaem/podoroot/domain/shared/inner"
	domain "github.com/podossaem/podoroot/domain/wallet"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podopaysql"
	"github.com/podossaem/podoroot/infra/entity"
)

type (
	walletRepository struct {
		client *podopaysql.Client
	}
)

func (r *walletRepository) InsertOne(
	ctx inner.Context,
	wallet domain.Wallet,
) (
	domain.Wallet,
	error,
) {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeWallet(wallet)

	result := db.Create(&e)

	if result.Error != nil {
		return domain.Wallet{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func NewWalletRepository(
	client *podopaysql.Client,
) domain.WalletRepository {
	return &walletRepository{
		client: client,
	}
}
