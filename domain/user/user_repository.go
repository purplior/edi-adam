package user

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	UserRepository interface {
		FindOneByID(
			ctx inner.Context,
			id string,
		) (
			User,
			error,
		)

		FindOneByAccount(
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
