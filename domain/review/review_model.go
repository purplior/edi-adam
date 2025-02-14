package review

import (
	"time"

	"github.com/purplior/sbec/domain/user"
)

type (
	Review struct {
		ID          string    `json:"id"`
		AuthorID    string    `json:"authorId"`
		AssistantID string    `json:"assistantId"`
		Content     string    `json:"content"`
		Score       float64   `json:"score"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		Author      user.User `json:"author"`
	}

	ReviewQueryOption struct {
		WithAuthor bool
	}

	ReviewInfo struct {
		ID        string             `json:"id"`
		Author    user.OtherUserInfo `json:"author"`
		Content   string             `json:"content"`
		Score     float64            `json:"score"`
		CreatedAt time.Time          `json:"createdAt"`
	}
)

func (r Review) ToInfo() ReviewInfo {
	return ReviewInfo{
		ID:        r.ID,
		Author:    r.Author.ToOtherUserInfo(),
		Content:   r.Content,
		Score:     r.Score,
		CreatedAt: r.CreatedAt,
	}
}

type (
	AddOneRequest struct {
		AuthorID    string  `json:"authorId"`
		AssistantID string  `json:"assistantId"`
		Content     string  `json:"content"`
		Score       float64 `json:"score"`
	}
)

func (r AddOneRequest) ToModelForInsert() Review {
	m := Review{
		AuthorID:    r.AuthorID,
		AssistantID: r.AssistantID,
		Content:     r.Content,
		Score:       r.Score,
	}

	return m
}
