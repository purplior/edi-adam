package assistant

import (
	"time"

	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/category"
	"github.com/purplior/podoroot/domain/user"
)

const (
	AssistantType_Writing AssistantType = 1
)

type (
	AssistantType int

	Assistant struct {
		ID                string              `json:"id"`
		ViewID            string              `json:"viewId"`
		AuthorID          string              `json:"authorId"`
		CategoryID        string              `json:"categoryId"`
		Assisters         []assister.Assister `json:"assisters"`
		AssistantType     AssistantType       `json:"assistantType"`
		Title             string              `json:"title"`
		Description       string              `json:"description"`
		IsPublic          bool                `json:"isPublic"`
		DefaultAssisterID string              `json:"defaultAssisterId"`
		CreatedAt         time.Time           `json:"createdAt"`
		Author            user.User           `json:"author"`
		Category          category.Category   `json:"category"`
	}

	AssistantJoinOption struct {
		WithAuthor    bool
		WithCategory  bool
		WithAssisters bool
	}
)

func (m *Assistant) ToInfo() (
	AssistantInfo,
	error,
) {
	return AssistantInfo{
		ViewID:            m.ViewID,
		Title:             m.Title,
		AssistantType:     m.AssistantType,
		Description:       m.Description,
		AuthorInfo:        m.Author.ToInfo(),
		CategoryInfo:      m.Category.ToInfo(),
		DefaultAssisterID: m.DefaultAssisterID,
		CreatedAt:         m.CreatedAt,
	}, nil
}

func (m *Assistant) ToDetail() (
	AssistantDetail,
	error,
) {
	assisterInfos := make([]assister.AssisterInfo, len(m.Assisters))
	for i, assister := range m.Assisters {
		assisterInfos[i] = assister.ToInfo()
	}

	return AssistantDetail{
		ViewID:            m.ViewID,
		AuthorInfo:        m.Author.ToInfo(),
		AssistantType:     m.AssistantType,
		Title:             m.Title,
		Description:       m.Description,
		IsPublic:          m.IsPublic,
		DefaultAssisterID: m.DefaultAssisterID,
		AssisterInfos:     assisterInfos,
		CreatedAt:         m.CreatedAt,
	}, nil
}

type (
	AssistantInfo struct {
		ViewID            string                `json:"viewId"`
		Title             string                `json:"title"`
		AssistantType     AssistantType         `json:"assistantType"`
		Description       string                `json:"description"`
		AuthorInfo        user.UserInfo         `json:"authorInfo"`
		CategoryInfo      category.CategoryInfo `json:"categoryInfo"`
		DefaultAssisterID string                `json:"defaultAssisterId"`
		CreatedAt         time.Time             `json:"createdAt"`
	}
)

type (
	AssistantDetail struct {
		ViewID            string                  `json:"viewId"`
		AssistantType     AssistantType           `json:"assistantType"`
		AuthorInfo        user.UserInfo           `json:"authorInfo"`
		Title             string                  `json:"title"`
		Description       string                  `json:"description"`
		IsPublic          bool                    `json:"isPublic"`
		AssisterInfos     []assister.AssisterInfo `json:"assisterInfos"`
		DefaultAssisterID string                  `json:"defaultAssisterId"`
		CreatedAt         time.Time               `json:"createdAt"`
	}
)

type (
	RegisterOneRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    bool   `json:"isPublic"`
	}
)
