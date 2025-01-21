package assistant

import (
	"time"

	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/category"
	"github.com/purplior/podoroot/domain/review"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/domain/user"
	"github.com/purplior/podoroot/lib/strgen"
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
		ID             string                    `json:"id"`
		ViewID         string                    `json:"viewId"`
		AuthorID       string                    `json:"authorId"`
		CategoryID     string                    `json:"categoryId"`
		AssisterID     string                    `json:"assisterId"`
		AssistantType  AssistantType             `json:"assistantType"`
		Icon           string                    `json:"icon"`
		Title          string                    `json:"title"`
		Description    string                    `json:"description"`
		Notice         string                    `json:"notice"`
		Tags           []string                  `json:"tags"`
		MetaTags       []string                  `json:"metaTags"`
		IsPublic       bool                      `json:"isPublic"`
		Status         AssistantStatus           `json:"status"`
		CreatedAt      time.Time                 `json:"createdAt"`
		UpdatedAt      time.Time                 `json:"updatedAt"`
		PublishedAt    *time.Time                `json:"publishedAt"`
		Author         user.User                 `json:"author"`
		Assister       assister.Assister         `json:"assister"`
		Category       category.Category         `json:"category"`
		Reviews        []review.Review           `json:"reviews"`
		ReviewPageMeta pagination.PaginationMeta `json:"reviewPageMeta"`
	}

	AssistantJoinOption struct {
		WithAuthor   bool
		WithCategory bool
		WithAssister bool
		WithReviews  bool
	}

	AssistantInfo struct {
		ID            string                `json:"id"`
		ViewID        string                `json:"viewId"`
		Icon          string                `json:"icon"`
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
		ID               string                    `json:"id"`
		ViewID           string                    `json:"viewId"`
		AssistantType    AssistantType             `json:"assistantType"`
		AuthorInfo       user.UserInfo             `json:"authorInfo"`
		Icon             string                    `json:"icon"`
		Title            string                    `json:"title"`
		Description      string                    `json:"description"`
		Notice           string                    `json:"notice"`
		Tags             []string                  `json:"tags"`
		MetaTags         []string                  `json:"metaTags"`
		IsPublic         bool                      `json:"isPublic"`
		IsMyRecentReview bool                      `json:"isMyRecentReview"`
		Status           AssistantStatus           `json:"status"`
		AssisterInfo     assister.AssisterInfo     `json:"assisterInfo"`
		ReviewInfos      []review.ReviewInfo       `json:"reviewInfos"`
		ReviewPageMeta   pagination.PaginationMeta `json:"reviewPageMeta"`
		CreatedAt        time.Time                 `json:"createdAt"`
		PublishedAt      *time.Time                `json:"publishedAt"`
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

	info := AssistantInfo{
		ID:            m.ID,
		ViewID:        m.ViewID,
		Icon:          m.Icon,
		Title:         m.Title,
		AssistantType: m.AssistantType,
		Description:   m.Description,
		Tags:          m.Tags,
		Status:        m.Status,
		AuthorInfo:    authorInfo,
		CategoryInfo:  categoryInfo,
		CreatedAt:     m.CreatedAt,
	}

	return info
}

func (m Assistant) ToDetail() AssistantDetail {
	authorInfo := m.Author.ToInfo()
	authorInfo.ID = m.AuthorID

	assisterInfo := m.Assister.ToInfo()
	assisterInfo.ID = m.AssisterID

	detail := AssistantDetail{
		ID:             m.ID,
		ViewID:         m.ViewID,
		AuthorInfo:     authorInfo,
		AssistantType:  m.AssistantType,
		Icon:           m.Icon,
		Title:          m.Title,
		Description:    m.Description,
		Notice:         m.Notice,
		Tags:           m.Tags,
		MetaTags:       m.MetaTags,
		IsPublic:       m.IsPublic,
		Status:         m.Status,
		AssisterInfo:   assisterInfo,
		ReviewPageMeta: m.ReviewPageMeta,
		CreatedAt:      m.CreatedAt,
		PublishedAt:    m.PublishedAt,
	}

	detail.ReviewInfos = make([]review.ReviewInfo, len(m.Reviews))
	for i, review := range m.Reviews {
		detail.ReviewInfos[i] = review.ToInfo()
	}

	return detail
}

type (
	RegisterOneRequest struct {
		Icon          string                          `json:"icon"`
		Title         string                          `json:"title"`
		Description   string                          `json:"description"`
		Notice        string                          `json:"notice"`
		CategoryID    string                          `json:"categoryId"`
		Tags          []string                        `json:"tags"`
		Fields        []assister.AssisterField        `json:"fields"`
		QueryMessages []assister.AssisterQueryMessage `json:"queryMessages"`
		Tests         []assister.AssisterInput        `json:"tests"`
		IsPublic      bool                            `json:"isPublic"`
	}
)

func (r RegisterOneRequest) ToModelForInsert(
	authorID string,
	assisterID string,
) Assistant {
	status := AssistantStatus_Registered
	if r.IsPublic {
		status = AssistantStatus_UnderReview
	}

	return Assistant{
		ViewID:        strgen.ShortUniqueID(),
		AuthorID:      authorID,
		CategoryID:    r.CategoryID,
		AssisterID:    assisterID,
		AssistantType: AssistantType_Formal,
		Icon:          r.Icon,
		Title:         r.Title,
		Description:   r.Description,
		Notice:        r.Notice,
		Tags:          r.Tags,
		MetaTags:      []string{},
		// 공개는 심사를 통해서만 수정됨.
		IsPublic: false,
		Status:   status,
	}
}

type (
	UpdateOneRequest struct {
		ID            string                          `json:"id"`
		Icon          string                          `json:"icon"`
		Title         string                          `json:"title"`
		Description   string                          `json:"description"`
		Notice        string                          `json:"notice"`
		CategoryID    string                          `json:"categoryId"`
		Tags          []string                        `json:"tags"`
		Fields        []assister.AssisterField        `json:"fields"`
		QueryMessages []assister.AssisterQueryMessage `json:"queryMessages"`
		Tests         []assister.AssisterInput        `json:"tests"`
		IsPublic      bool                            `json:"isPublic"`
	}
)

func (r UpdateOneRequest) ToModelForUpdate(
	existedAssistant Assistant,
) Assistant {
	existedAssistant.Icon = r.Icon
	existedAssistant.CategoryID = r.CategoryID
	existedAssistant.Title = r.Title
	existedAssistant.Description = r.Description
	existedAssistant.Notice = r.Notice
	existedAssistant.Tags = r.Tags
	if !existedAssistant.IsPublic && r.IsPublic {
		existedAssistant.Status = AssistantStatus_UnderReview
	}

	return existedAssistant
}
