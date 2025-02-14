package repository

import (
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	domain "github.com/purplior/sbec/domain/wallet"
	"github.com/purplior/sbec/infra/database"
	"github.com/purplior/sbec/infra/database/sqldb"
	"github.com/purplior/sbec/infra/entity"
)

type (
	walletRepository struct {
		client *sqldb.Client
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
	client *sqldb.Client,
) domain.WalletRepository {
	return &walletRepository{
		client: client,
	}
}
