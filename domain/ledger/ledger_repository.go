package ledger

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	LedgerRepository interface {
		InsertOne(
			ctx inner.Context,
			ledger Ledger,
		) (
			Ledger,
			error,
		)
	}
)
