package verification

import "github.com/podossaem/root/domain/context"

type (
	EmailVerificationRepository interface {
		Create(
			ctx context.APIContext,
			emailVerification EmailVerification,
		) (
			EmailVerification,
			error,
		)
	}
)
