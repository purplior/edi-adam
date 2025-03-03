package repository

import (
	"github.com/purplior/edi-adam/domain/shared/model"
	domain "github.com/purplior/edi-adam/domain/walletlog"
	"github.com/purplior/edi-adam/infra/database/dynamo"
)

type (
	walletLogRepository struct {
		dynamoRepository[model.WalletLog, domain.QueryOption]
	}
)

func NewWalletLogRepository(
	client *dynamo.Client,
) domain.WalletLogRepository {
	var repo dynamoRepository[model.WalletLog, domain.QueryOption]

	repo.client = client

	return &walletLogRepository{
		dynamoRepository: repo,
	}
}
