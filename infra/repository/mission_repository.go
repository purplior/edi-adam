package repository

import (
	domain "github.com/purplior/edi-adam/domain/mission"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	missionRepository struct {
		postgreRepository[model.Mission, domain.QueryOption]
	}
)

func NewMissionRepository(
	client *postgre.Client,
) domain.MissionRepository {
	var repo postgreRepository[model.Mission, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		return db, nil
	}

	return &missionRepository{
		postgreRepository: repo,
	}
}
