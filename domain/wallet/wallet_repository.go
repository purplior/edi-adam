package wallet

import (
	"github.com/purplior/podoroot/domain/shared/inner"
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

		FindOne_ByUserID(
			ctx inner.Context,
			userID string,
		) (
			Wallet,
			error,
		)

		UpdateOne_ByUserIDAndDelta(
			ctx inner.Context,
			userID string,
			podoDelta int,
		) (
			Wallet,
			error,
		)
	}
)
