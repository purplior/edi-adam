package assistant

import (
	"time"

	"github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/domain/user"
)

type (
	Assistant struct {
		ID          string              `json:"id"`
		AuthorID    string              `json:"authorId"`
		Author      user.User           `json:"author"`
		Assisters   []assister.Assister `json:"assisters"`
		Title       string              `json:"title"`
		Description string              `json:"description"`
		IsPublic    bool                `json:"isPublic"`
		CreatedAt   time.Time           `json:"createdAt"`
	}

	AssistantJoinOption struct {
		WithAuthor   bool
		WithAssister bool
	}
)

func (m *Assistant) ToInfo() (
	AssistantInfo,
	error,
) {
	return AssistantInfo{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		AuthorInfo:  m.Author.ToInfo(),
		CreatedAt:   m.CreatedAt,
	}, nil
}

func (m *Assistant) ToDetail() (
	AssistantDetail,
	error,
) {
	return AssistantDetail{
		ID:          m.ID,
		AuthorInfo:  m.Author.ToInfo(),
		Title:       m.Title,
		Description: m.Description,
		IsPublic:    m.IsPublic,
		CreatedAt:   m.CreatedAt,
	}, nil
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
	AssistantDetail struct {
		ID          string        `json:"id"`
		AuthorInfo  user.UserInfo `json:"authorInfo"`
		Title       string        `json:"title"`
		Description string        `json:"description"`
		IsPublic    bool          `json:"isPublic"`
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
