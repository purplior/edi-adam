package assistant

import (
	"time"

	"github.com/podossaem/podoroot/domain/user"
)

type (
	Assistant struct {
		ID          string    `json:"id"`
		AuthorID    string    `json:"authorId"`
		Author      user.User `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		IsPublic    bool      `json:"isPublic"`
		CreatedAt   time.Time `json:"createdAt"`
	}
)

func (m *Assistant) ToInfo() AssistantInfo {
	return AssistantInfo{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		AuthorInfo:  m.Author.ToInfo(),
		CreatedAt:   m.CreatedAt,
	}
}

type (
	AssistantInfo struct {
		ID          string        `json:"id"`
		Title       string        `json:"title"`
		Description string        `json:"description"`
		AuthorInfo  user.UserInfo `json:"authorInfo"`
		CreatedAt   time.Time     `json:"createdAt"`
	}
)

type (
	RegisterOneRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    bool   `json:"isPublic"`
	}
)
