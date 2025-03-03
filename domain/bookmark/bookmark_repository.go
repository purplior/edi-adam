package bookmark

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	BookmarkRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Bookmark,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Bookmark,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.Bookmark,
		) (
			model.Bookmark,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.Bookmark,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
