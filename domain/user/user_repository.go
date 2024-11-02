package user

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	UserRepository interface {
		FindOne_ByID(
			ctx inner.Context,
			id string,
		) (
			User,
			error,
		)

		FindOne_ByAccount(
			ctx inner.Context,
			joinMethod string,
			accountID string,
		) (
			User,
			error,
		)

		InsertOne(
			ctx inner.Context,
			user User,
		) (
			User,
			error,
		)
	}
)
