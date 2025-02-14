package user

import (
	"strings"

	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/lib/mydate"
)

type (
	UserService interface {
		GetOne_ByID(
			ctx inner.Context,
			id string,
		) (
			User,
			error,
		)

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

		CheckNicknameExistence(
			ctx inner.Context,
			nickname string,
		) (
			bool,
			error,
		)

		Inactive(
			ctx inner.Context,
			userID string,
		) error

		UpdateOne_Password_ByAccount(
			ctx inner.Context,
			joinMethod string,
			accountID string,
			newPassword string,
		) error
	}

	userService struct {
		userRepository UserRepository
	}
)

func (s *userService) GetOne_ByID(
	ctx inner.Context,
	id string,
) (
	User,
	error,
) {
	return s.userRepository.FindOne_ByID(ctx, id)
}

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

func (s *userService) CheckNicknameExistence(
	ctx inner.Context,
	nickname string,
) (
	bool,
	error,
) {
	if strings.Contains(nickname, "포도쌤") {
		return false, exception.ErrNotAllowedNickname
	}

	_, err := s.userRepository.FindOne_ByNickname(
		ctx,
		nickname,
	)
	if err != nil {
		if err == exception.ErrNoRecord {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *userService) Inactive(
	ctx inner.Context,
	userID string,
) error {
	return s.userRepository.UpdateOne_InactivatedFields(
		ctx,
		userID,
		true,
		mydate.Now(),
	)
}

func (s *userService) UpdateOne_Password_ByAccount(
	ctx inner.Context,
	joinMethod string,
	accountID string,
	newPassword string,
) error {
	return s.userRepository.UpdateOne_Password_ByAccount(
		ctx,
		joinMethod,
		accountID,
		newPassword,
	)
}

func NewUserService(
	userRepository UserRepository,
) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
