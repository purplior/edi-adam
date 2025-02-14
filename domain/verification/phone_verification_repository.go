package verification

import "github.com/purplior/sbec/domain/shared/inner"

type (
	PhoneVerificationRepository interface {
		InsertOne(
			ctx inner.Context,
			phoneVerification PhoneVerification,
		) (
			PhoneVerification,
			error,
		)

		FindOneById(
			ctx inner.Context,
			id string,
		) (
			PhoneVerification,
			error,
		)

		FindRecentOne_ByPhoneNumber(
			ctx inner.Context,
			phoneNumber string,
		) (
			PhoneVerification,
			error,
		)

		FindCount_ByPhoneNumber(
			ctx inner.Context,
			phoneNumber string,
		) (
			count int,
			err error,
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
