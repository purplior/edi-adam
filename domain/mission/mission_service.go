package mission

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	MissionService interface {
		GetPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Mission,
			pagination.PaginationMeta,
			error,
		)
	}
)

type (
	missionService struct {
		missionRepository MissionRepository
	}
)

func (s *missionService) GetPaginatedList(
	session inner.Session,
	query pagination.PaginationQuery[QueryOption],
) (
	[]model.Mission,
	pagination.PaginationMeta,
	error,
) {
	return s.missionRepository.ReadPaginatedList(
		session,
		query,
	)
}

func NewMissionService(
	missionRepository MissionRepository,
) MissionService {
	return &missionService{
		missionRepository: missionRepository,
	}
}
