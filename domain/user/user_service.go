package user

import (
	"github.com/podossaem/podoroot/domain/shared/inner"
)

type (
	UserService interface {
		GetOne_ByAccount(
			ctx inner.Context,
			joinMethod string,
			accountID string,
		) (
			User,
			error,
		)

		GetDetailOne_ByID(
			ctx inner.Context,
			id string,
		) (
			UserDetail,
			error,
		)

		RegisterOne(
			ctx inner.Context,
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

func (s *userService) GetOne_ByAccount(
	ctx inner.Context,
	joinMethod string,
	accountID string,
) (
	User,
	error,
) {
	return s.userRepository.FindOne_ByAccount(
		ctx,
		joinMethod,
		accountID,
	)
}

func (s *userService) GetDetailOne_ByID(
	ctx inner.Context,
	id string,
) (
	UserDetail,
	error,
) {
	user, err := s.userRepository.FindOne_ByID(ctx, id)
	if err != nil {
		return UserDetail{}, err
	}

	return user.ToDetail(), nil
}

func (s *userService) RegisterOne(
	ctx inner.Context,
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
