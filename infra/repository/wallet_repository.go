package repository

import (
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	domain "github.com/purplior/podoroot/domain/wallet"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
)

type (
	walletRepository struct {
		client *podosql.Client
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

func (r *walletRepository) FindOne_ByUserID(
	ctx inner.Context,
	userID string,
) (
	domain.Wallet,
	error,
) {
	db := r.client.DBWithContext(ctx)

	var wallet entity.Wallet
	if err := db.Where("owner_id = ?", userID).First(&wallet).Error; err != nil {
		return domain.Wallet{}, database.ToDomainError(err)
	}

	return wallet.ToModel(), nil
}

func (r *walletRepository) UpdateOne_ByUserIDAndDelta(
	ctx inner.Context,
	userID string,
	podoDelta int,
) (
	domain.Wallet,
	error,
) {
	db := r.client.DBWithContext(ctx)

	var wallet entity.Wallet
	if err := db.Where("owner_id = ?", userID).First(&wallet).Error; err != nil {
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
	client *podosql.Client,
) domain.WalletRepository {
	return &walletRepository{
		client: client,
	}
}
