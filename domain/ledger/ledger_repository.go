package ledger

import "github.com/purplior/podoroot/domain/shared/inner"

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
