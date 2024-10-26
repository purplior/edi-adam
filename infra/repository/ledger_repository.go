package repository

import (
	domain "github.com/podossaem/podoroot/domain/ledger"
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podopaysql"
	"github.com/podossaem/podoroot/infra/entity"
)

type (
	ledgerRepository struct {
		client *podopaysql.Client
	}
)

func (r *ledgerRepository) InsertOne(
	ctx context.APIContext,
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
	client *podopaysql.Client,
) domain.LedgerRepository {
	return &ledgerRepository{
		client: client,
	}
}
