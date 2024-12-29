package verification

import "github.com/purplior/podoroot/domain/shared/inner"

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

		FindRecentOneByPhoneNumber(
			ctx inner.Context,
			phoneNumber string,
		) (
			PhoneVerification,
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
