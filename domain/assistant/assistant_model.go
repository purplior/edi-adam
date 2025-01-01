package assistant

import (
	"time"

	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/category"
	"github.com/purplior/podoroot/domain/user"
)

const (
	AssistantType_Formal AssistantType = 1

	AssistantStatus_Registered  AssistantStatus = "registered"
	AssistantStatus_UnderReview AssistantStatus = "under_review"
	AssistantStatus_Approved    AssistantStatus = "approved"
	AssistantStatus_Rejected    AssistantStatus = "rejected"
)

type (
	AssistantType   int
	AssistantStatus string

	Assistant struct {
		ID            string            `json:"id"`
		ViewID        string            `json:"viewId"`
		AuthorID      string            `json:"authorId"`
		CategoryID    string            `json:"categoryId"`
		AssisterID    string            `json:"assisterId"`
		AssistantType AssistantType     `json:"assistantType"`
		Title         string            `json:"title"`
		Description   string            `json:"description"`
		Tags          []string          `json:"tags"`
		MetaTags      []string          `json:"metaTags"`
		IsPublic      bool              `json:"isPublic"`
		Status        AssistantStatus   `json:"status"`
		CreatedAt     time.Time         `json:"createdAt"`
		UpdatedAt     time.Time         `json:"updatedAt"`
		Author        user.User         `json:"author"`
		Assister      assister.Assister `json:"assister"`
		Category      category.Category `json:"category"`
	}

	AssistantJoinOption struct {
		WithAuthor   bool
		WithCategory bool
		WithAssister bool
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
		ID            string                `json:"id"`
		ViewID        string                `json:"viewId"`
		AssistantType AssistantType         `json:"assistantType"`
		AuthorInfo    user.UserInfo         `json:"authorInfo"`
		Title         string                `json:"title"`
		Description   string                `json:"description"`
		Tags          []string              `json:"tags"`
		MetaTags      []string              `json:"metaTags"`
		IsPublic      bool                  `json:"isPublic"`
		Status        AssistantStatus       `json:"status"`
		AssisterInfo  assister.AssisterInfo `json:"assisterInfo"`
		CreatedAt     time.Time             `json:"createdAt"`
	}
)

func (m Assistant) Copy() Assistant {
	return m
}

func (m Assistant) ToInfo() AssistantInfo {
	authorInfo := m.Author.ToInfo()
	authorInfo.ID = m.AuthorID

	categoryInfo := m.Category.ToInfo()
	categoryInfo.ID = m.CategoryID

	return AssistantInfo{
		ID:            m.ID,
		ViewID:        m.ViewID,
		Title:         m.Title,
		AssistantType: m.AssistantType,
		Description:   m.Description,
		Tags:          m.Tags,
		Status:        m.Status,
		AuthorInfo:    authorInfo,
		CategoryInfo:  categoryInfo,
		CreatedAt:     m.CreatedAt,
	}
}

func (m Assistant) ToDetail() AssistantDetail {
	authorInfo := m.Author.ToInfo()
	authorInfo.ID = m.AuthorID

	assisterInfo := m.Assister.ToInfo()
	assisterInfo.ID = m.AssisterID

	return AssistantDetail{
		ID:            m.ID,
		ViewID:        m.ViewID,
		AuthorInfo:    authorInfo,
		AssistantType: m.AssistantType,
		Title:         m.Title,
		Description:   m.Description,
		Tags:          m.Tags,
		MetaTags:      m.MetaTags,
		IsPublic:      m.IsPublic,
		Status:        m.Status,
		AssisterInfo:  assisterInfo,
		CreatedAt:     m.CreatedAt,
	}
}

type (
	RegisterOneRequest struct {
		Title         string                          `json:"title"`
		Description   string                          `json:"description"`
		CategoryID    string                          `json:"categoryId"`
		Tags          []string                        `json:"tags"`
		Fields        []assister.AssisterField        `json:"fields"`
		QueryMessages []assister.AssisterQueryMessage `json:"queryMessages"`
		Tests         []assister.AssisterInput        `json:"tests"`
		IsPublic      bool                            `json:"isPublic"`
	}

	UpdateOneRequest struct {
		ID            string                          `json:"id"`
		Title         string                          `json:"title"`
		Description   string                          `json:"description"`
		CategoryID    string                          `json:"categoryId"`
		Tags          []string                        `json:"tags"`
		Fields        []assister.AssisterField        `json:"fields"`
		QueryMessages []assister.AssisterQueryMessage `json:"queryMessages"`
		Tests         []assister.AssisterInput        `json:"tests"`
		IsPublic      bool                            `json:"isPublic"`
	}
)
