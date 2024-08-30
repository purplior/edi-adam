package verification

import "github.com/podossaem/root/domain/context"

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
	return EmailVerification{}, nil
}

func NewEmailVerificationService(
	repository EmailVerificationRepository,
) EmailVerificationService {
	return &service{
		repository: repository,
	}
}
