package user

import (
	"time"

	"github.com/purplior/podoroot/domain/shared/exception"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID              string    `json:"id"`
		JoinMethod      string    `json:"joinMethod"`
		AccountID       string    `json:"accountId"`
		AccountPassword string    `json:"accountPassword"`
		AvatarTheme     int       `json:"avatarTheme"`
		AvatarText      string    `json:"avatarText"`
		Nickname        string    `json:"nickname"`
		Role            int       `json:"role"`
		CreatedAt       time.Time `json:"createdAt"`
	}

	UserInfo struct {
		ID          string `json:"id"`
		AvatarTheme int    `json:"avatarTheme"`
		AvatarText  string `json:"avatarText"`
		Nickname    string `json:"nickname"`
	}

	UserDetail struct {
		ID          string    `json:"id"`
		JoinMethod  string    `json:"joinMethod"`
		AccountID   string    `json:"accountId"`
		AvatarTheme int       `json:"avatarTheme"`
		AvatarText  string    `json:"avatarText"`
		Nickname    string    `json:"nickname"`
		CreatedAt   time.Time `json:"createdAt"`
	}
)

func (m User) ToInfo() UserInfo {
	return UserInfo{
		ID:          m.ID,
		AvatarTheme: m.AvatarTheme,
		AvatarText:  m.AvatarText,
		Nickname:    m.Nickname,
	}
}

func (m User) ToDetail() UserDetail {
	return UserDetail{
		ID:          m.ID,
		JoinMethod:  m.JoinMethod,
		AccountID:   m.AccountID,
		AvatarTheme: m.AvatarTheme,
		AvatarText:  m.AvatarText,
		Nickname:    m.Nickname,
		CreatedAt:   m.CreatedAt,
	}
}

func (e User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(e.AccountPassword), []byte(password)); err != nil {
		return exception.ErrUnauthorized
	}
	return nil
}

func (e *User) HashPassword() error {
	cost := 12
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(e.AccountPassword), cost)
	if err != nil {
		return err
	}

	e.AccountPassword = string(hashedPassword)

	return nil
}

const (
	JoinMethod_Email = "email"

	Role_User   = 100
	Role_Master = 10000

	ID_Podo = "1"
)
