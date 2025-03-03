package repository

import (
	domain "github.com/purplior/edi-adam/domain/assister"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/dynamo"
)

type (
	assisterRepository struct {
		dynamoRepository[model.Assister, domain.QueryOption]
	}
)

func NewAssisterRepository(
	client *dynamo.Client,
) domain.AssisterRepository {
	var repo dynamoRepository[model.Assister, domain.QueryOption]

	repo.client = client

	return &assisterRepository{
		dynamoRepository: repo,
	}
}
