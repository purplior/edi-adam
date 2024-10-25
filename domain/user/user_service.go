package user

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	UserService interface {
		GetOneByAccount(
			ctx context.APIContext,
			joinMethod string,
			accountID string,
		) (
			User,
			error,
		)

		GetDetailOneByID(
			ctx context.APIContext,
			id string,
		) (
			UserDetail,
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

func (s *userService) GetOneByAccount(
	ctx context.APIContext,
	joinMethod string,
	accountID string,
) (
	User,
	error,
) {
	return s.userRepository.FindOneByAccount(
		ctx,
		joinMethod,
		accountID,
	)
}

func (s *userService) GetDetailOneByID(
	ctx context.APIContext,
	id string,
) (
	UserDetail,
	error,
) {
	user, err := s.userRepository.FindOneByID(ctx, id)
	if err != nil {
		return UserDetail{}, err
	}

	return user.ToDetail(), nil
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
