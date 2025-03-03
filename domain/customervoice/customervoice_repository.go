package customervoice

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	CustomerVoiceRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.CustomerVoice,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.CustomerVoice,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.CustomerVoice,
		) (
			model.CustomerVoice,
			error,
		)
	}
)
