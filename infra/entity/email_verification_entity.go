package entity

import (
	"time"

	"github.com/podossaem/podoroot/domain/verification"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	EmailVerification struct {
		ID         uint      `gorm:"primaryKey;autoIncrement"`
		Email      string    `gorm:"size:100;not null"`
		Code       string    `gorm:"size:10;not null"`
		IsVerified bool      `gorm:"not null"`
		IsConsumed bool      `gorm:"not null"`
		ExpiredAt  time.Time `gorm:"not null"`
		CreatedAt  time.Time `gorm:"autoCreateTime"`
	}
)

func (e EmailVerification) ToModel() verification.EmailVerification {
	model := verification.EmailVerification{
		Email:      e.Email,
		Code:       e.Code,
		IsVerified: e.IsVerified,
		IsConsumed: e.IsConsumed,
		ExpiredAt:  e.ExpiredAt,
		CreatedAt:  e.CreatedAt,
	}

	if e.ID > 0 {
		model.ID = dt.Str(e.ID)
	}

	return model
}

func MakeEmailVerification(m verification.EmailVerification) EmailVerification {
	entity := EmailVerification{
		Email:      m.Email,
		Code:       m.Code,
		IsVerified: m.IsVerified,
		IsConsumed: m.IsConsumed,
		ExpiredAt:  m.ExpiredAt,
		CreatedAt:  m.CreatedAt,
	}

	if len(m.ID) > 0 {
		entity.ID = dt.UInt(m.ID)
	}

	return entity
}
