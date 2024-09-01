package verification

import "github.com/podossaem/podoroot/domain/context"

type (
	EmailVerificationRepository interface {
		InsertOne(
			ctx context.APIContext,
			emailVerification EmailVerification,
		) (
			EmailVerification,
			error,
		)

		FindOneById(
			ctx context.APIContext,
			id string,
		) (
			EmailVerification,
			error,
		)

		FindOneByEmail(
			ctx context.APIContext,
			email string,
		) (
			EmailVerification,
			error,
		)

		UpdateOne_IsVerified(
			ctx context.APIContext,
			id string,
			isVerified bool,
		) (
			EmailVerification,
			error,
		)

		UpdateOne_isConsumed(
			ctx context.APIContext,
			id string,
			isConsumed bool,
		) (
			EmailVerification,
			error,
		)
	}
)
