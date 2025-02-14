package repository

import (
	domain "github.com/purplior/sbec/domain/ledger"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/infra/database"
	"github.com/purplior/sbec/infra/database/sqldb"
	"github.com/purplior/sbec/infra/entity"
)

type (
	ledgerRepository struct {
		client *sqldb.Client
	}
)

func (r *ledgerRepository) InsertOne(
	ctx inner.Context,
	ledger domain.Ledger,
) (
	domain.Ledger,
	error,
) {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeLedger(ledger)

	result := db.Create(&e)

	if result.Error != nil {
		return domain.Ledger{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func NewLedgerRepository(
	client *sqldb.Client,
) domain.LedgerRepository {
	return &ledgerRepository{
		client: client,
	}
}
