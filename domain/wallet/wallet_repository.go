package wallet

import (
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	WalletRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Wallet,
			error,
		)

		Create(
			session inner.Session,
			m model.Wallet,
		) (
			model.Wallet,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.Wallet,
		) error

		UpdatesCoinDelta_ByOwnerID(
			session inner.Session,
			userID uint,
			delta int,
		) (
			model.Wallet,
			error,
		)

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
