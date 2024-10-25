package user

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	UserRepository interface {
		FindOneByID(
			ctx context.APIContext,
			id string,
		) (
			User,
			error,
		)

		FindOneByAccount(
			ctx context.APIContext,
			joinMethod string,
			accountID string,
		) (
			User,
			error,
		)

		InsertOne(
			ctx context.APIContext,
			user User,
		) (
			User,
			error,
		)
	}
)
