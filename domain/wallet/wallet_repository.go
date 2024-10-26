package wallet

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	WalletRepository interface {
		InsertOne(
			ctx context.APIContext,
			wallet Wallet,
		) (
			Wallet,
			error,
		)
	}
)
