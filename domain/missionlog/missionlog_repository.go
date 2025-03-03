package missionlog

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	MissionLogRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.MissionLog,
			error,
		)

		ReadList(
			session inner.Session,
			queryOption QueryOption,
		) (
			[]model.MissionLog,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.MissionLog,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.MissionLog,
		) (
			model.MissionLog,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.MissionLog,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
