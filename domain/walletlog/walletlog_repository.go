package walletlog

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	WalletLogRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.WalletLog,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.WalletLog,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.WalletLog,
		) (
			model.WalletLog,
			error,
		)
	}
)
