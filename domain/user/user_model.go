package user

import (
	"time"

	"github.com/purplior/podoroot/domain/shared/exception"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID               string     `json:"id"`
		JoinMethod       string     `json:"joinMethod"`
		AccountID        string     `json:"accountId"`
		AccountPassword  string     `json:"accountPassword"`
		AvatarTheme      int        `json:"avatarTheme"`
		AvatarText       string     `json:"avatarText"`
		Nickname         string     `json:"nickname"`
		PhoneNumber      string     `json:"phoneNumber"`
		Role             int        `json:"role"`
		IsMarketingAgree bool       `json:"isMarketingAgree"`
		IsInactivated    bool       `json:"isInactivated"`
		CreatedAt        time.Time  `json:"createdAt"`
		InactivatedAt    *time.Time `json:"inactivatedAt"`
	}

	UserInfo struct {
		ID          string `json:"id"`
		AvatarTheme int    `json:"avatarTheme"`
		AvatarText  string `json:"avatarText"`
		Nickname    string `json:"nickname"`
	}

	OtherUserInfo struct {
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
		PhoneNumber string    `json:"phoneNumber"`
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

func (m User) ToOtherUserInfo() OtherUserInfo {
	return OtherUserInfo{
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
		PhoneNumber: m.PhoneNumber,
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
	JoinMethod_Email       = "email"
	JoinMethod_PhoneNumber = "phone"

	Role_User   = 100
	Role_Master = 10000

	ID_Podo = "1"
)
