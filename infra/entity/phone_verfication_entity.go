package entity

import (
	"time"

	"github.com/purplior/sbec/domain/verification"
	"github.com/purplior/sbec/lib/dt"
)

type (
	PhoneVerification struct {
		ID          uint      `gorm:"primaryKey;autoIncrement"`
		PhoneNumber string    `gorm:"size:20;not null"`
		Code        string    `gorm:"size:10;not null"`
		ReferenceID string    `gorm:"size:80"`
		IsVerified  bool      `gorm:"not null"`
		IsConsumed  bool      `gorm:"not null"`
		ExpiredAt   time.Time `gorm:"not null"`
		CreatedAt   time.Time `gorm:"autoCreateTime"`
	}
)

func (e PhoneVerification) ToModel() verification.PhoneVerification {
	model := verification.PhoneVerification{
		PhoneNumber: e.PhoneNumber,
		Code:        e.Code,
		ReferenceID: e.ReferenceID,
		IsVerified:  e.IsVerified,
		IsConsumed:  e.IsConsumed,
		ExpiredAt:   e.ExpiredAt,
		CreatedAt:   e.CreatedAt,
	}

	if e.ID > 0 {
		model.ID = dt.Str(e.ID)
	}

	return model
}

func MakePhoneVerification(m verification.PhoneVerification) PhoneVerification {
	entity := PhoneVerification{
		PhoneNumber: m.PhoneNumber,
		Code:        m.Code,
		ReferenceID: m.ReferenceID,
		IsVerified:  m.IsVerified,
		IsConsumed:  m.IsConsumed,
		ExpiredAt:   m.ExpiredAt,
		CreatedAt:   m.CreatedAt,
	}

	if len(m.ID) > 0 {
		entity.ID = dt.UInt(m.ID)
	}

	return entity
}
