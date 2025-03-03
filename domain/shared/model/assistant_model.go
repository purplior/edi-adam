package model

import (
	"time"
)

type (
	// @ENUM {샘비서 상태}
	AssistantStatus int
)

const (
	AssistantStatus_Registered  AssistantStatus = 1
	AssistantStatus_UnderReview AssistantStatus = 2
	AssistantStatus_Approved    AssistantStatus = 3
	AssistantStatus_Rejected    AssistantStatus = 4
)

type (
	// @MODEL {샘비서 모델 : 기본 모델}
	Assistant struct {
		ID                uint            `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		Icon              string          `gorm:"size:80" json:"icon"`
		Title             string          `gorm:"size:80" json:"title"`
		Description       string          `gorm:"size:255" json:"description"`
		Notice            string          `gorm:"size:255" json:"notice"`
		Tags              []string        `gorm:"serializer:json" json:"tags"`
		Status            AssistantStatus `gorm:"size:20" json:"status"`
		IsPublic          bool            `json:"isPublic"`
		CurrentAssisterID string          `json:"currentAssisterId,omitempty"`

		CreatedAt   time.Time  `json:"createdAt"`
		UpdatedAt   time.Time  `json:"updatedAt"`
		PublishedAt *time.Time `json:"publishedAt"`

		AuthorID   uint      `json:"authorId,omitempty"`
		Author     *User     `json:"author,omitempty"`
		CategoryID string    `gorm:"size:20" json:"categoryId,omitempty"`
		Category   *Category `json:"category,omitempty"`
	}

	// @MODEL {샘비서 모델 : 카드}
	AssistantCard struct {
		ID          uint            `json:"id"`
		Icon        string          `json:"icon"`
		Title       string          `json:"title"`
		Description string          `json:"description"`
		Tags        []string        `json:"tags"`
		Status      AssistantStatus `json:"status"`

		CreatedAt time.Time `json:"createdAt"`

		Author   UserProfile  `json:"authorProfile"`
		Category CategoryChip `json:"category"`
	}
)

// @METHOD {카드 정보로 변환}
func (m Assistant) ToCard() AssistantCard {
	authorProfile := m.Author.ToProfile()
	authorProfile.ID = m.AuthorID

	categoryChip := m.Category.ToChip()
	categoryChip.ID = m.CategoryID

	card := AssistantCard{
		ID:          m.ID,
		Icon:        m.Icon,
		Title:       m.Title,
		Description: m.Description,
		Tags:        m.Tags,
		Status:      m.Status,

		CreatedAt: m.CreatedAt,

		Author:   authorProfile,
		Category: categoryChip,
	}

	return card
}
