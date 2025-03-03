package assistant

import "github.com/purplior/edi-adam/lib/validator"

type (
	// DTO {샘비서 조회 옵션}
	QueryOption struct {
		ID         uint
		AuthorID   uint
		CategoryID string
		IsPublic   bool

		WithAuthor   bool
		WithCategory bool
	}

	// @DTO {샘비서 등록 요청}
	RegisterDTO struct {
		UserID            uint
		CurrentAssisterID string

		Icon        string
		Title       string
		Description string
		Notice      string
		CategoryID  string
		Tags        []string
		IsPublic    bool
	}

	// @DTO {샘비서 수정 요청}
	UpdateDTO struct {
		ID     uint
		UserID uint

		Icon        string
		Title       string
		Description string
		Notice      string
		CategoryID  string
		Tags        []string
		IsPublic    bool
	}

	RemoveDTO struct {
		ID     uint
		UserID uint
	}

	ApproveDTO struct {
		ID uint
	}
)

func (m RegisterDTO) IsValid() bool {
	if !validator.CheckValidAssistantTitle(m.Title) {
		return false
	}
	if !validator.CheckValidAssistantDescription(m.Description) {
		return false
	}
	if !validator.CheckValidAssistantTags(m.Tags) {
		return false
	}

	return true
}
