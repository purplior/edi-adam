package repository

import (
	domain "github.com/purplior/edi-adam/domain/customervoice"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	customerVoiceRepository struct {
		postgreRepository[model.CustomerVoice, domain.QueryOption]
	}
)

func NewCustomerVoiceRepository(
	client *postgre.Client,
) domain.CustomerVoiceRepository {
	var repo postgreRepository[model.CustomerVoice, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		return db, nil
	}

	return &customerVoiceRepository{
		postgreRepository: repo,
	}
}
