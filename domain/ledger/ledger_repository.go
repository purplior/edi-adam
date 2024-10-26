package ledger

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	LedgerRepository interface {
		InsertOne(
			ctx context.APIContext,
			ledger Ledger,
		) (
			Ledger,
			error,
		)
	}
)
