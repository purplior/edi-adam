package entity

import (
	"time"

	domain "github.com/purplior/podoroot/domain/assistant"
	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	Assistant struct {
		ID            uint   `gorm:"primaryKey;autoIncrement"`
		ViewID        string `gorm:"size:36;not null;unique"`
		AuthorID      uint
		CategoryID    uint
		AssistantType uint      `gorm:"default:0"`
		Title         string    `gorm:"size:80;not null"`  // 20자 이내
		Description   string    `gorm:"size:255;not null"` // 80자 이내
		Tags          []string  `gorm:"serializer:json"`
		MetaTags      []string  `gorm:"serializer:json"`
		IsPublic      bool      `gorm:"default:false;not null"`
		Status        string    `gorm:"size:80"`
		CreatedAt     time.Time `gorm:"autoCreateTime"`
		Author        User
		Category      Category
		Assisters     []Assister `gorm:"foreignKey:AssistantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Bookmarks     []Bookmark `gorm:"foreignKey:AssistantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (e Assistant) ToModel() domain.Assistant {
	model := domain.Assistant{
		ViewID:        e.ViewID,
		AssistantType: domain.AssistantType(e.AssistantType),
		Title:         e.Title,
		Description:   e.Description,
		Tags:          e.Tags,
		MetaTags:      e.MetaTags,
		IsPublic:      e.IsPublic,
		Status:        domain.AssistantStatus(e.Status),
		CreatedAt:     e.CreatedAt,
	}

	if e.ID > 0 {
		model.ID = dt.Str(e.ID)
	}
	if e.AuthorID > 0 {
		model.AuthorID = dt.Str(e.AuthorID)
		model.Author = e.Author.ToModel()
	}
	if e.CategoryID > 0 {
		model.CategoryID = dt.Str(e.CategoryID)
		model.Category = e.Category.ToModel()
	}

	assisters := make([]assister.Assister, len(e.Assisters))
	for i, e := range e.Assisters {
		assisters[i] = e.ToModel()
	}
	model.Assisters = assisters

	return model
}

func MakeAssistant(m domain.Assistant) Assistant {
	entity := Assistant{
		ViewID:        m.ViewID,
		AssistantType: uint(m.AssistantType),
		Title:         m.Title,
		Description:   m.Description,
		Tags:          m.Tags,
		MetaTags:      m.MetaTags,
		IsPublic:      m.IsPublic,
		Status:        string(m.Status),
		CreatedAt:     m.CreatedAt,
	}

	if len(m.ID) > 0 {
		entity.ID = dt.UInt(m.ID)
	}
	if len(m.AuthorID) > 0 {
		entity.AuthorID = dt.UInt(m.AuthorID)
	}
	if len(m.CategoryID) > 0 {
		entity.CategoryID = dt.UInt(m.CategoryID)
	}

	return entity
}
