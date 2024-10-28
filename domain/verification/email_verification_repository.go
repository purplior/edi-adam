package verification

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	EmailVerificationRepository interface {
		InsertOne(
			ctx inner.Context,
			emailVerification EmailVerification,
		) (
			EmailVerification,
			error,
		)

		FindOneById(
			ctx inner.Context,
			id string,
		) (
			EmailVerification,
			error,
		)

		FindRecentOneByEmail(
			ctx inner.Context,
			email string,
		) (
			EmailVerification,
			error,
		)

		UpdateOne_IsVerified(
			ctx inner.Context,
			id string,
			isVerified bool,
		) error

		UpdateOne_isConsumed(
			ctx inner.Context,
			id string,
			isConsumed bool,
		) error
	}
)
