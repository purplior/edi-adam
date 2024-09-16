package persist

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/verification"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	EmailVerification struct {
		ID         uint      `gorm:"primaryKey;autoIncrement"`
		Email      string    `gorm:"type:varchar(100);not null"`
		Code       string    `gorm:"type:varchar(6);not null"`
		IsVerified bool      `gorm:"not null"`
		IsConsumed bool      `gorm:"not null"`
		ExpiredAt  time.Time `gorm:"not null"`
		CreatedAt  time.Time `gorm:"autoCreateTime"`
	}
)

func (e EmailVerification) ToModel() domain.EmailVerification {
	model := domain.EmailVerification{
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

func MakeEmailVerification(m domain.EmailVerification) EmailVerification {
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
