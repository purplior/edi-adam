package persist

import (
	"github.com/podossaem/root/domain/context"
	"github.com/podossaem/root/domain/verification"
)

type (
	repository struct{}
)

func (r *repository) Create(
	ctx context.APIContext,
	emailVerification verification.EmailVerification,
) (
	verification.EmailVerification,
	error,
) {
	return verification.EmailVerification{}, nil
}

func NewEmailVerificationRepository() verification.EmailVerificationRepository {
	return &repository{}
}
