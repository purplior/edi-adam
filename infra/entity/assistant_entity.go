package entity

import (
	"time"

	"github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Assistant struct {
		ID                uint   `gorm:"primaryKey;autoIncrement"`
		ViewID            string `gorm:"type:varchar(36);not null;unique"`
		AuthorID          uint
		Author            User
		Assisters         []Assister `gorm:"foreignKey:AssistantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Title             string     `gorm:"type:varchar(80);not null"`  // 20자 이내
		Description       string     `gorm:"type:varchar(255);not null"` // 80자 이내
		IsPublic          bool       `gorm:"default:false;not null"`
		DefaultAssisterID uint
		CreatedAt         time.Time `gorm:"autoCreateTime"`
	}
)

func (e Assistant) ToModel() assistant.Assistant {
	model := assistant.Assistant{
		ViewID:      e.ViewID,
		Title:       e.Title,
		Description: e.Description,
		IsPublic:    e.IsPublic,
		CreatedAt:   e.CreatedAt,
	}

	if e.ID > 0 {
		model.ID = dt.Str(e.ID)
	}
	if e.AuthorID > 0 {
		model.AuthorID = dt.Str(e.AuthorID)
		model.Author = e.Author.ToModel()
	}

	assisters := make([]assister.Assister, len(e.Assisters))
	for i, e := range e.Assisters {
		assisters[i] = e.ToModel()
	}
	model.Assisters = assisters

	if e.DefaultAssisterID > 0 {
		model.DefaultAssisterID = dt.Str(e.DefaultAssisterID)
	}

	return model
}

func MakeAssistant(m assistant.Assistant) Assistant {
	entity := Assistant{
		Title:       m.Title,
		Description: m.Description,
		IsPublic:    m.IsPublic,
		CreatedAt:   m.CreatedAt,
	}

	if len(m.ID) > 0 {
		entity.ID = dt.UInt(m.ID)
	}
	if len(m.AuthorID) > 0 {
		entity.AuthorID = dt.UInt(m.AuthorID)
	}

	return entity
}
