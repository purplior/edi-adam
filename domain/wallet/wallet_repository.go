package wallet

import (
	"github.com/podossaem/podoroot/domain/shared/inner"
)

type (
	WalletRepository interface {
		InsertOne(
			ctx inner.Context,
			wallet Wallet,
		) (
			Wallet,
			error,
		)

		UpdateOneByUserIDAndDelta(
			ctx inner.Context,
			userId string,
			podoDelta int,
		) (
			Wallet,
			error,
		)
	}
)
