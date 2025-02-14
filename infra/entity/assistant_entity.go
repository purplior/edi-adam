package entity

import (
	"time"

	domain "github.com/purplior/sbec/domain/assistant"
	"github.com/purplior/sbec/domain/review"
	"github.com/purplior/sbec/lib/dt"
)

type (
	Assistant struct {
		ID            uint   `gorm:"primaryKey;autoIncrement"`
		ViewID        string `gorm:"size:36;not null;unique"`
		AuthorID      uint
		CategoryID    uint
		AssisterID    string    `gorm:"size:80"`
		AssistantType uint      `gorm:"default:0"`
		Icon          string    `gorm:"size:80"`
		Title         string    `gorm:"size:80;not null"`  // 20자 이내
		Description   string    `gorm:"size:255;not null"` // 80자 이내
		Notice        string    `gorm:"size:255"`
		Tags          []string  `gorm:"serializer:json"`
		MetaTags      []string  `gorm:"serializer:json"`
		IsPublic      bool      `gorm:"default:false;not null"`
		Status        string    `gorm:"size:80"`
		CreatedAt     time.Time `gorm:"autoCreateTime"`
		UpdatedAt     time.Time `gorm:"autoUpdateTime"`
		PublishedAt   *time.Time
		Author        User
		Category      Category
		Bookmarks     []Bookmark `gorm:"foreignKey:AssistantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Reviews       []Review   `gorm:"foreignKey:AssistantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (e Assistant) ToModel() domain.Assistant {
	model := domain.Assistant{
		ViewID:        e.ViewID,
		AssisterID:    e.AssisterID,
		AssistantType: domain.AssistantType(e.AssistantType),
		Icon:          e.Icon,
		Title:         e.Title,
		Description:   e.Description,
		Notice:        e.Notice,
		Tags:          e.Tags,
		MetaTags:      e.MetaTags,
		IsPublic:      e.IsPublic,
		Status:        domain.AssistantStatus(e.Status),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
		PublishedAt:   e.PublishedAt,
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

	model.Reviews = make([]review.Review, len(e.Reviews))
	for i, eReview := range e.Reviews {
		model.Reviews[i] = eReview.ToModel()
	}

	return model
}

func MakeAssistant(m domain.Assistant) Assistant {
	entity := Assistant{
		ViewID:        m.ViewID,
		AssisterID:    m.AssisterID,
		AssistantType: uint(m.AssistantType),
		Icon:          m.Icon,
		Title:         m.Title,
		Description:   m.Description,
		Notice:        m.Notice,
		Tags:          m.Tags,
		MetaTags:      m.MetaTags,
		IsPublic:      m.IsPublic,
		Status:        string(m.Status),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		PublishedAt:   m.PublishedAt,
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

	entity.Reviews = make([]Review, len(m.Reviews))
	for i, review := range m.Reviews {
		entity.Reviews[i] = MakeReview(review)
	}

	return entity
}
