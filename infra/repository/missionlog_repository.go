package repository

import (
	domain "github.com/purplior/edi-adam/domain/missionlog"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	missionLogRepository struct {
		postgreRepository[model.MissionLog, domain.QueryOption]
	}
)

func NewMissionLogRepository(
	client *postgre.Client,
) domain.MissionLogRepository {
	var repo postgreRepository[model.MissionLog, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		query := db
		if opt.WithMission {
			query = query.Preload("Mission")
		}
		if opt.ID > 0 {
			query = query.Where("id = ?", opt.ID)
		}
		if opt.MissionID > 0 {
			query = query.Where("mission_id = ?", opt.MissionID)
		}
		if opt.UserID > 0 {
			query = query.Where("user_id = ?", opt.UserID)
		}

		return query, nil
	}

	return &missionLogRepository{
		postgreRepository: repo,
	}
}
