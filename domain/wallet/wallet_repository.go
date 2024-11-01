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

		FindOneByUserID(
			ctx inner.Context,
			userID string,
		) (
			Wallet,
			error,
		)

		UpdateOneByUserIDAndDelta(
			ctx inner.Context,
			userID string,
			podoDelta int,
		) (
			Wallet,
			error,
		)
	}
)
