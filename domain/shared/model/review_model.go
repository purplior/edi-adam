package model

import "time"

type (
	Review struct {
		ID      uint    `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		Content string  `gorm:"size:1500" json:"content"`
		Score   float64 `json:"score"`

		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
		UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

		Author      *User `json:"author,omitempty"`
		AuthorID    uint  `json:"authorId,omitempty"`
		AssistantID uint  `json:"assistantId,omitempty"`
	}
)

type (
	ReviewCard struct {
		ID      uint    `json:"id"`
		Content string  `json:"content"`
		Score   float64 `json:"score"`

		CreatedAt time.Time `json:"createdAt"`

		Author UserProfile `json:"author"`
	}
)

func (r Review) ToCard() ReviewCard {
	return ReviewCard{
		ID:        r.ID,
		Content:   r.Content,
		Score:     r.Score,
		CreatedAt: r.CreatedAt,
		Author:    r.Author.ToProfile(),
	}
}
