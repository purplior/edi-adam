package repository

import (
	"github.com/podossaem/podoroot/domain/shared/exception"
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

func (r *walletRepository) UpdateOneByUserIDAndDelta(
	ctx inner.Context,
	userId string,
	podoDelta int,
) (
	domain.Wallet,
	error,
) {
	db := r.client.DBWithContext(ctx)

	var wallet entity.Wallet
	if err := db.Where("owner_id = ?", userId).First(&wallet).Error; err != nil {
		return domain.Wallet{}, database.ToDomainError(err)
	}

	if wallet.Podo+int64(podoDelta) < 0 {
		return domain.Wallet{}, exception.ErrNoPodo
	}

	result := db.Model(&wallet).Update("podo", wallet.Podo+int64(podoDelta))
	if result.Error != nil {
		return domain.Wallet{}, database.ToDomainError(result.Error)
	}

	return wallet.ToModel(), nil
}

func NewWalletRepository(
	client *podopaysql.Client,
) domain.WalletRepository {
	return &walletRepository{
		client: client,
	}
}
