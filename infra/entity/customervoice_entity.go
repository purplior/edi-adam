package entity

import (
	"time"

	domain "github.com/purplior/podoroot/domain/customervoice"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	CustomerVoice struct {
		ID        uint `gorm:"primaryKey;autoIncrement"`
		UserID    uint
		Title     string    `gorm:"size:255"`
		Content   string    `gorm:"size:1000"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}
)

func (e CustomerVoice) ToModel() domain.CustomerVoice {
	model := domain.CustomerVoice{
		Title:     e.Title,
		Content:   e.Content,
		CreatedAt: e.CreatedAt,
	}

	if e.ID > 0 {
		model.ID = dt.Str(e.ID)
	}
	if e.UserID > 0 {
		model.UserID = dt.Str(e.UserID)
	}

	return model
}

func MakeCustomerVoice(m domain.CustomerVoice) CustomerVoice {
	entity := CustomerVoice{
		Title:     m.Title,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
	}

	if len(m.ID) > 0 {
		entity.ID = dt.UInt(m.ID)
	}
	if len(m.UserID) > 0 {
		entity.UserID = dt.UInt(m.UserID)
	}

	return entity
}
