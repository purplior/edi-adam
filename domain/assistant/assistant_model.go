package assistant

import (
	"time"

	"github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/domain/user"
)

type (
	Assistant struct {
		ID              string              `json:"id"`
		ViewID          string              `json:"viewId"`
		AuthorID        string              `json:"authorId"`
		Author          user.User           `json:"author"`
		Assisters       []assister.Assister `json:"assisters"`
		Title           string              `json:"title"`
		Description     string              `json:"description"`
		IsPublic        bool                `json:"isPublic"`
		DefaultAssister assister.Assister   `json:"defaultAssister"`
		CreatedAt       time.Time           `json:"createdAt"`
	}

	AssistantJoinOption struct {
		WithAuthor          bool
		WithAssisters       bool
		WithDefaultAssister bool
	}
)

func (m *Assistant) ToInfo() (
	AssistantInfo,
	error,
) {
	return AssistantInfo{
		ViewID:      m.ViewID,
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
		ViewID:              m.ViewID,
		AuthorInfo:          m.Author.ToInfo(),
		Title:               m.Title,
		Description:         m.Description,
		IsPublic:            m.IsPublic,
		DefaultAssisterInfo: m.DefaultAssister.ToInfo(),
		CreatedAt:           m.CreatedAt,
	}, nil
}

type (
	AssistantInfo struct {
		ViewID      string        `json:"viewId"`
		Title       string        `json:"title"`
		Description string        `json:"description"`
		AuthorInfo  user.UserInfo `json:"authorInfo"`
		CreatedAt   time.Time     `json:"createdAt"`
	}
)

type (
	AssistantDetail struct {
		ViewID              string                `json:"viewId"`
		AuthorInfo          user.UserInfo         `json:"authorInfo"`
		Title               string                `json:"title"`
		Description         string                `json:"description"`
		IsPublic            bool                  `json:"isPublic"`
		DefaultAssisterInfo assister.AssisterInfo `json:"defaultAssisterInfo"`
		CreatedAt           time.Time             `json:"createdAt"`
	}
)

type (
	RegisterOneRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    bool   `json:"isPublic"`
	}
)
