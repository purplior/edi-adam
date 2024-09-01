package user

import (
	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/domain/verification"
)

type (
	UserService interface {
		SignUpByEmailVerification(
			ctx context.APIContext,
			request SignUpRequest,
		) (
			User,
			error,
		)
	}

	userService struct {
		emailVerificationService verification.EmailVerificationService
		userRepository           UserRepository
	}
)

func (s *userService) SignUpByEmailVerification(
	ctx context.APIContext,
	request SignUpRequest,
) (
	User,
	error,
) {
	verification, err := s.emailVerificationService.Consume(
		ctx,
		request.VerificationID,
	)
	if err != nil {
		return User{}, err
	}

	return s.userRepository.InsertOne(
		ctx,
		User{
			JoinMethod:      JoinMethod_Email,
			AccountID:       verification.Email,
			AccountPassword: request.Password,
			Nickname:        request.Nickname,
			Role:            Role_User,
		},
	)
}

func NewUserService(
	emailVerificationService verification.EmailVerificationService,
	userRepository UserRepository,
) UserService {
	return &userService{
		emailVerificationService: emailVerificationService,
		userRepository:           userRepository,
	}
}
