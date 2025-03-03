package assistant

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	AssistantRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Assistant,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Assistant,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.Assistant,
		) (
			model.Assistant,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.Assistant,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
