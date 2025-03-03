package repository

import (
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	domain "github.com/purplior/edi-adam/domain/wallet"
	"github.com/purplior/edi-adam/infra/database"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	walletRepository struct {
		postgreRepository[model.Wallet, domain.QueryOption]
	}
)

func (r *walletRepository) UpdatesCoinDelta_ByOwnerID(
	session inner.Session,
	ownerID uint,
	delta int,
) (
	m model.Wallet,
	err error,
) {
	db := r.client.DBWithContext(session)

	if m, err = r.Read(
		session,
		domain.QueryOption{OwnerID: ownerID},
	); err != nil {
		return m, err
	}

	newCoin := m.Coin + int64(delta)
	if newCoin < 0 {
		return m, exception.ErrNoCoin
	}

	result := db.Model(&m).Update("podo", newCoin)
	if result.Error != nil {
		return m, database.ToDomainError(result.Error)
	}

	return m, nil
}

func NewWalletRepository(
	client *postgre.Client,
) domain.WalletRepository {
	var repo postgreRepository[model.Wallet, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		return db, nil
	}

	return &walletRepository{
		postgreRepository: repo,
	}
}
