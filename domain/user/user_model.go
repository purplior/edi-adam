package user

import (
	"time"

	"github.com/podossaem/podoroot/domain/exception"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID              string    `json:"id"`
		JoinMethod      string    `json:"joinMethod"`
		AccountID       string    `json:"accountId"`
		AccountPassword string    `json:"accountPassword"`
		Nickname        string    `json:"nickname"`
		Role            int       `json:"role"`
		CreatedAt       time.Time `json:"createdAt"`
	}

	UserInfo struct {
		ID       string `json:"id"`
		Nickname string `json:"nickname"`
	}
)

func (m *User) ToInfo() UserInfo {
	return UserInfo{
		ID:       m.ID,
		Nickname: m.Nickname,
	}
}

func (e *User) ComparePassword(password string) error {
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
