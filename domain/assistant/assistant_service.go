package assistant

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/lib/mydate"
)

type (
	AssistantService interface {
		Get(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Assistant,
			error,
		)

		GetPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Assistant,
			pagination.PaginationMeta,
			error,
		)

		Register(
			session inner.Session,
			dto RegisterDTO,
		) (
			model.Assistant,
			error,
		)

		Update(
			session inner.Session,
			dto UpdateDTO,
		) error

		Remove(
			session inner.Session,
			dto RemoveDTO,
		) error

		Approve(
			session inner.Session,
			dto ApproveDTO,
		) error
	}
)

type (
	assistantService struct {
		assistantRepository AssistantRepository
	}
)

func (s *assistantService) Get(
	session inner.Session,
	queryOption QueryOption,
) (
	model.Assistant,
	error,
) {
	return s.assistantRepository.Read(
		session,
		queryOption,
	)
}

func (s *assistantService) GetPaginatedList(
	session inner.Session,
	query pagination.PaginationQuery[QueryOption],
) (
	[]model.Assistant,
	pagination.PaginationMeta,
	error,
) {
	return s.assistantRepository.ReadPaginatedList(
		session,
		query,
	)
}

func (s *assistantService) Register(
	session inner.Session,
	dto RegisterDTO,
) (
	m model.Assistant,
	err error,
) {
	status := model.AssistantStatus_Registered
	if dto.IsPublic {
		status = model.AssistantStatus_UnderReview
	}

	m, err = s.assistantRepository.Create(
		session,
		model.Assistant{
			Icon:        dto.Icon,
			Title:       dto.Title,
			Description: dto.Description,
			Notice:      dto.Notice,
			Tags:        dto.Tags,
			Status:      status,
			IsPublic:    false,

			AuthorID:          dto.UserID,
			CategoryID:        dto.CategoryID,
			CurrentAssisterID: dto.CurrentAssisterID,
		},
	)

	return m, err
}

func (s *assistantService) Update(
	session inner.Session,
	dto UpdateDTO,
) (err error) {
	var m model.Assistant
	if m, err = s.assistantRepository.Read(
		session,
		QueryOption{
			ID: dto.ID,
		},
	); err != nil {
		return err
	}
	if dto.UserID != m.AuthorID {
		return exception.ErrBadRequest
	}

	m.Icon = dto.Icon
	m.Title = dto.Title
	m.Description = dto.Description
	m.Notice = dto.Notice
	m.CategoryID = dto.CategoryID
	m.Tags = dto.Tags

	return s.assistantRepository.Updates(
		session,
		QueryOption{ID: dto.ID},
		m,
	)
}

func (s *assistantService) Remove(
	session inner.Session,
	dto RemoveDTO,
) (err error) {
	var m model.Assistant
	if m, err = s.assistantRepository.Read(
		session,
		QueryOption{ID: dto.ID},
	); err != nil {
		return err
	}

	if m.AuthorID != dto.UserID {
		return exception.ErrBadRequest
	}
	// 공개된 상태의 샘비서는 지울 수 없음
	if m.IsPublic {
		return exception.ErrBadRequest
	}

	if err := s.assistantRepository.Delete(
		session,
		QueryOption{ID: dto.ID},
	); err != nil {
		return err
	}

	return nil
}

func (s *assistantService) Approve(
	session inner.Session,
	dto ApproveDTO,
) (err error) {
	var m model.Assistant
	if m, err = s.assistantRepository.Read(
		session,
		QueryOption{ID: dto.ID},
	); err != nil {
		return err
	}

	// 리뷰 중 상태일 경우만 approve 가능함
	if m.Status != model.AssistantStatus_UnderReview {
		return exception.ErrBadRequest
	}

	now := mydate.Now()
	m.Status = model.AssistantStatus_Approved
	m.IsPublic = true
	m.PublishedAt = &now

	return s.assistantRepository.Updates(
		session,
		QueryOption{ID: dto.ID},
		m,
	)
}

func NewAssistantService(
	assistantRepository AssistantRepository,
) AssistantService {
	return &assistantService{
		assistantRepository: assistantRepository,
	}
}
