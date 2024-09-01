package user

import "github.com/podossaem/podoroot/domain/context"

type (
	UserRepository interface {
		InsertOne(
			ctx context.APIContext,
			user User,
		) (
			User,
			error,
		)
	}
)
