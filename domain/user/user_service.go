package user

import "github.com/podossaem/podoroot/domain/context"

type (
	UserService interface {
		GetByAccount(
			ctx context.APIContext,
			joinMethod string,
			accountID string,
		) (
			User,
			error,
		)

		RegisterOne(
			ctx context.APIContext,
			user User,
		) (
			newUser User,
			err error,
		)
	}

	userService struct {
		userRepository UserRepository
	}
)

func (s *userService) GetByAccount(
	ctx context.APIContext,
	joinMethod string,
	accountID string,
) (
	User,
	error,
) {
	return s.userRepository.FindByAccount(
		ctx,
		joinMethod,
		accountID,
	)
}

func (s *userService) RegisterOne(
	ctx context.APIContext,
	user User,
) (
	newUser User,
	err error,
) {
	if err := user.HashPassword(); err != nil {
		return User{}, err
	}

	return s.userRepository.InsertOne(ctx, user)
}

func NewUserService(
	userRepository UserRepository,
) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
