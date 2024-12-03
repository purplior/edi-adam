package repository

import (
	domain "github.com/purplior/podoroot/domain/ledger"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
)

type (
	ledgerRepository struct {
		client *podosql.Client
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
	client *podosql.Client,
) domain.LedgerRepository {
	return &ledgerRepository{
		client: client,
	}
}
