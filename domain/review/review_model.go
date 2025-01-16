package review

import (
	"time"
)

type (
	Review struct {
		ID          string    `json:"id"`
		AuthorID    string    `json:"authorId"`
		AssistantID string    `json:"assistantId"`
		Content     string    `json:"content"`
		Score       float64   `json:"score"`
		CreatedAt   time.Time `json:"createdAt"`
	}
)

type (
	ReviewJoinOption struct {
		WithAuthor bool
	}
)

type (
	WriteOneRequest struct {
		AuthorID    string  `json:"authorId"`
		AssistantID string  `json:"assistantId"`
		Content     string  `json:"content"`
		Score       float64 `json:"score"`
	}
)

func (r WriteOneRequest) ToModelForInsert() Review {
	m := Review{
		AuthorID:    r.AuthorID,
		AssistantID: r.AssistantID,
		Content:     r.Content,
		Score:       r.Score,
	}

	return m
}
