package repository

import (
	domain "github.com/purplior/podoroot/domain/customervoice"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
)

type (
	customerVoiceRepository struct {
		client *podosql.Client
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
	client *podosql.Client,
) domain.CustomerVoiceRepository {
	return &customerVoiceRepository{
		client: client,
	}
}
