package repository

import (
	domain "github.com/podossaem/podoroot/domain/ledger"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
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
