package verification

import (
	"github.com/podossaem/root/domain/context"
	"github.com/podossaem/root/lib/strgen"
)

type (
	EmailVerificationService interface {
		RequestCode(
			ctx context.APIContext,
			email string,
		) (
			EmailVerification,
			error,
		)
	}

	service struct {
		repository EmailVerificationRepository
	}
)

func (s *service) RequestCode(
	ctx context.APIContext,
	email string,
) (EmailVerification, error) {
	verification := EmailVerification{
		Email:      email,
		Code:       strgen.RandomNumber(6),
		IsConsumed: false,
		IsVerified: false,
	}

	return s.repository.Create(ctx, verification)
}

func NewEmailVerificationService(
	repository EmailVerificationRepository,
) EmailVerificationService {
	return &service{
		repository: repository,
	}
}
