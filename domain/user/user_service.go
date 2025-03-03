package user

import (
	"strings"

	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/lib/mydate"
)

type (
	UserService interface {
		Get(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.User,
			error,
		)

		RegisterMember(
			session inner.Session,
			request RegisterMemberDTO,
		) (
			newUser model.User,
			err error,
		)

		CheckNicknameExistence(
			session inner.Session,
			nickname string,
		) (
			exist bool,
			err error,
		)

		Inactive(
			session inner.Session,
			userID uint,
		) (
			err error,
		)
	}

	userService struct {
		userRepository UserRepository
	}
)

func (s *userService) Get(
	session inner.Session,
	queryOption QueryOption,
) (
	m model.User,
	err error,
) {
	accountID, err := PhoneNumberToAccountID(queryOption.PhoneNumber)
	if err != nil {
		return m, err
	}

	queryOption.AccountID = accountID
	queryOption.PhoneNumber = ""
	m, err = s.userRepository.Read(session, queryOption)

	return m, err
}

func (s *userService) RegisterMember(
	session inner.Session,
	dto RegisterMemberDTO,
) (
	m model.User,
	err error,
) {
	accountID, err := PhoneNumberToAccountID(dto.PhoneNumber)
	if err != nil {
		return model.User{}, err
	}

	m, err = s.userRepository.Create(session, model.User{
		AccountID:        accountID,
		Nickname:         dto.Nickname,
		Avatar:           dto.Avatar,
		Role:             model.UserRole_Member,
		IsMarketingAgree: dto.IsMarketingAgree,
	})

	return m, err
}

func (s *userService) CheckNicknameExistence(
	session inner.Session,
	nickname string,
) (
	bool,
	error,
) {
	if strings.Contains(nickname, "샘비서") {
		return false, exception.ErrNotAllowedNickname
	}

	_, err := s.userRepository.Read(
		session,
		QueryOption{Nickname: nickname},
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
	session inner.Session,
	userID uint,
) error {
	now := mydate.Now()
	return s.userRepository.Updates(
		session,
		QueryOption{
			ID: userID,
		},
		model.User{
			IsInactivated: true,
			InactivatedAt: &now,
		},
	)
}

func NewUserService(
	userRepository UserRepository,
) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
