package ledger

import "github.com/purplior/sbec/domain/shared/inner"

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
