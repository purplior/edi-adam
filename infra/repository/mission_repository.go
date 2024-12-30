package repository

import (
	domain "github.com/purplior/podoroot/domain/mission"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
)

type (
	missionRepository struct {
		client *podosql.Client
	}
)

func (r *missionRepository) FindOne_ByIDAndUserID(
	ctx inner.Context,
	id string,
	userID string,
) (
	domain.Mission,
	error,
) {
	var e entity.Mission
	db := r.client.DBWithContext(ctx)

	if err := db.Where("id = ?", id).First(&e).Error; err != nil {
		return domain.Mission{}, err
	}

	if err := db.Preload("Challenges", "user_id = ?", userID).
		Where("missions.id = ?", id).
		First(&e).Error; err != nil {
		return domain.Mission{}, database.ToDomainError(err)
	}

	if len(e.Challenges) == 0 {
		return domain.Mission{}, exception.ErrNoRecord
	}

	return e.ToModel(), nil
}

func (r *missionRepository) FindPaginatedList_OnlyPublic_ByUserID(
	ctx inner.Context,
	userID string,
	page int,
	pageSize int,
) (
	[]domain.Mission,
	pagination.PaginationMeta,
	error,
) {
	var entities []entity.Mission
	var totalCount int64
	db := r.client.DBWithContext(ctx)

	if err := db.Model(&entity.Mission{}).
		Where("is_public = ?", true).
		Count(&totalCount).Error; err != nil {
		return nil, pagination.PaginationMeta{}, database.ToDomainError(err)
	}

	offset := (page - 1) * pageSize
	if err := db.
		Offset(offset).
		Limit(pageSize).
		Where("is_public", true).
		Find(&entities).Error; err != nil {
		return nil, pagination.PaginationMeta{}, database.ToDomainError(err)
	}

	meta := pagination.PaginationMeta{
		Page:      page,
		Size:      pageSize,
		TotalPage: int((totalCount + int64(pageSize) - 1) / int64(pageSize)),
	}

	missions := make([]domain.Mission, len(entities))
	for i, entity := range entities {
		missions[i] = entity.ToModel()
	}

	return missions, meta, nil
}

func NewMissionRepository(
	client *podosql.Client,
) domain.MissionRepository {
	return &missionRepository{
		client: client,
	}
}
