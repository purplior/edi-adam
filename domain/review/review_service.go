package review

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/lib/mydate"
)

type (
	ReviewService interface {
		GetPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Review,
			pagination.PaginationMeta,
			error,
		)

		Add(
			session inner.Session,
			dto AddDTO,
		) (
			model.Review,
			error,
		)
	}
)

type (
	reviewService struct {
		reviewRepository ReviewRepository
	}
)

func (s *reviewService) GetPaginatedList(
	session inner.Session,
	query pagination.PaginationQuery[QueryOption],
) (
	mArr []model.Review,
	meta pagination.PaginationMeta,
	err error,
) {
	return s.reviewRepository.ReadPaginatedList(
		session,
		query,
	)
}

func (s *reviewService) Add(
	session inner.Session,
	dto AddDTO,
) (
	m model.Review,
	err error,
) {
	if m, err = s.reviewRepository.Read(
		session,
		QueryOption{
			AuthorID:    session.Identity().ID,
			AssistantID: dto.AssistantID,
		},
	); err != nil {
		if err == exception.ErrNoRecord {
			return s.addAndGetWithAuthor(session, dto)
		}

		return m, err
	}

	now := mydate.Now()
	prev := m.CreatedAt
	dayDiff := mydate.DaysDifference(now, prev)
	if dayDiff < 30 {
		return m, exception.ErrBadRequest
	}

	return s.addAndGetWithAuthor(session, dto)
}

func (s *reviewService) addAndGetWithAuthor(
	session inner.Session,
	dto AddDTO,
) (
	m model.Review,
	err error,
) {
	if m, err = s.reviewRepository.Create(
		session,
		model.Review{
			Content:     dto.Content,
			Score:       dto.Score,
			AuthorID:    session.Identity().ID,
			AssistantID: dto.AssistantID,
		},
	); err != nil {
		return m, err
	}

	if m, err = s.reviewRepository.Read(
		session,
		QueryOption{
			ID:         m.ID,
			WithAuthor: true,
		},
	); err != nil {
		return m, err
	}

	return m, nil
}

func NewReviewService(
	reviewRepository ReviewRepository,
) ReviewService {
	return &reviewService{
		reviewRepository: reviewRepository,
	}
}
