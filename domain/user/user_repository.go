package user

import (
	"time"

	"github.com/purplior/podoroot/domain/shared/inner"
)

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

		FindOne_ByNickname(
			ctx inner.Context,
			nickname string,
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

		UpdateOne_InactivatedFields(
			ctx inner.Context,
			userID string,
			isInactivated bool,
			inactivatedAt time.Time,
		) error
	}
)
