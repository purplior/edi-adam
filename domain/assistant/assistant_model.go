package assistant

import (
	"time"

	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/assisterform"
	"github.com/purplior/podoroot/domain/category"
	"github.com/purplior/podoroot/domain/user"
)

const (
	AssistantType_Formal AssistantType = 1

	AssistantStatus_Registered  AssistantStatus = "registered"
	AssistantStatus_UnderReview AssistantStatus = "under_review"
	AssistantStatus_Approved    AssistantStatus = "approved"
)

type (
	AssistantType   int
	AssistantStatus string

	Assistant struct {
		ID            string              `json:"id"`
		ViewID        string              `json:"viewId"`
		AuthorID      string              `json:"authorId"`
		CategoryID    string              `json:"categoryId"`
		Assisters     []assister.Assister `json:"assisters"`
		AssistantType AssistantType       `json:"assistantType"`
		Title         string              `json:"title"`
		Description   string              `json:"description"`
		Tags          []string            `json:"tags"`
		IsPublic      bool                `json:"isPublic"`
		Status        AssistantStatus     `json:"status"`
		CreatedAt     time.Time           `json:"createdAt"`
		Author        user.User           `json:"author"`
		Category      category.Category   `json:"category"`
	}

	AssistantJoinOption struct {
		WithAuthor    bool
		WithCategory  bool
		WithAssisters bool
	}

	AssistantInfo struct {
		ID            string                `json:"id"`
		ViewID        string                `json:"viewId"`
		Title         string                `json:"title"`
		AssistantType AssistantType         `json:"assistantType"`
		Description   string                `json:"description"`
		Tags          []string              `json:"tags"`
		Status        AssistantStatus       `json:"status"`
		AuthorInfo    user.UserInfo         `json:"authorInfo"`
		CategoryInfo  category.CategoryInfo `json:"categoryInfo"`
		CreatedAt     time.Time             `json:"createdAt"`
	}

	AssistantDetail struct {
		ID            string                  `json:"id"`
		ViewID        string                  `json:"viewId"`
		AssistantType AssistantType           `json:"assistantType"`
		AuthorInfo    user.UserInfo           `json:"authorInfo"`
		Title         string                  `json:"title"`
		Description   string                  `json:"description"`
		Tags          []string                `json:"tags"`
		IsPublic      bool                    `json:"isPublic"`
		Status        AssistantStatus         `json:"status"`
		AssisterInfos []assister.AssisterInfo `json:"assisterInfos"`
		CreatedAt     time.Time               `json:"createdAt"`
	}
)

func (m *Assistant) ToInfo() AssistantInfo {
	return AssistantInfo{
		ID:            m.ID,
		ViewID:        m.ViewID,
		Title:         m.Title,
		AssistantType: m.AssistantType,
		Description:   m.Description,
		Tags:          m.Tags,
		Status:        m.Status,
		AuthorInfo:    m.Author.ToInfo(),
		CategoryInfo:  m.Category.ToInfo(),
		CreatedAt:     m.CreatedAt,
	}
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
		ID:            m.ID,
		ViewID:        m.ViewID,
		AuthorInfo:    m.Author.ToInfo(),
		AssistantType: m.AssistantType,
		Title:         m.Title,
		Description:   m.Description,
		Tags:          m.Tags,
		IsPublic:      m.IsPublic,
		Status:        m.Status,
		AssisterInfos: assisterInfos,
		CreatedAt:     m.CreatedAt,
	}, nil
}

type (
	RegisterOneRequest struct {
		Title         string                              `json:"title"`
		Description   string                              `json:"description"`
		CategoryID    string                              `json:"categoryId"`
		Tags          []string                            `json:"tags"`
		Fields        []assisterform.AssisterField        `json:"fields"`
		QueryMessages []assisterform.AssisterQueryMessage `json:"queryMessages"`
		Tests         []assisterform.AssisterInput        `json:"tests"`
		IsPublic      bool                                `json:"isPublic"`
	}

	UpdateOneRequest struct {
		ID            string                              `json:"id"`
		Title         string                              `json:"title"`
		Description   string                              `json:"description"`
		CategoryID    string                              `json:"categoryId"`
		Tags          []string                            `json:"tags"`
		Fields        []assisterform.AssisterField        `json:"fields"`
		QueryMessages []assisterform.AssisterQueryMessage `json:"queryMessages"`
		Tests         []assisterform.AssisterInput        `json:"tests"`
		IsPublic      bool                                `json:"isPublic"`
	}
)
