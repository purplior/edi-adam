package user

import "github.com/podossaem/podoroot/domain/context"

type (
	UserRepository interface {
		FindByAccount(
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
