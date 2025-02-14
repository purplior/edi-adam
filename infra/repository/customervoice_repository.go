package repository

import (
	domain "github.com/purplior/sbec/domain/customervoice"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/infra/database"
	"github.com/purplior/sbec/infra/database/sqldb"
	"github.com/purplior/sbec/infra/entity"
)

type (
	customerVoiceRepository struct {
		client *sqldb.Client
	}
)

func (r *customerVoiceRepository) InsertOne(
	ctx inner.Context,
	customerVoice domain.CustomerVoice,
) (
	domain.CustomerVoice,
	error,
) {
	e := entity.MakeCustomerVoice(customerVoice)
	result := r.client.DBWithContext(ctx).Create(&e)

	if result.Error != nil {
		return domain.CustomerVoice{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func NewCustomerVoiceRepository(
	client *sqldb.Client,
) domain.CustomerVoiceRepository {
	return &customerVoiceRepository{
		client: client,
	}
}
