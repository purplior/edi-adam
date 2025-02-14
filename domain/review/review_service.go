package review

import (
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/shared/pagination"
	"github.com/purplior/sbec/lib/mydate"
)

type (
	ReviewService interface {
		GetPaginatedList_ByAssistantID(
			ctx inner.Context,
			assistantID string,
			pageRequest pagination.PaginationRequest,
		) (
			[]Review,
			pagination.PaginationMeta,
			error,
		)

		AddOne(
			ctx inner.Context,
			request AddOneRequest,
		) (
			Review,
			error,
		)

		UpdateRecentOne(
			ctx inner.Context,
			request AddOneRequest,
		) (
			Review,
			error,
		)
	}
)

type (
	reviewService struct {
		reviewRepository ReviewRepository
	}
)

func (s *reviewService) GetPaginatedList_ByAssistantID(
	ctx inner.Context,
	assistantID string,
	pageRequest pagination.PaginationRequest,
) (
	[]Review,
	pagination.PaginationMeta,
	error,
) {
	return s.reviewRepository.FindPaginatedList_ByAssistantID(
		ctx,
		assistantID,
		pageRequest,
	)
}

func (s *reviewService) AddOne(
	ctx inner.Context,
	request AddOneRequest,
) (
	Review,
	error,
) {
	existedReview, err := s.reviewRepository.FindOne_ByAuthorAndAssistantID(
		ctx,
		request.AuthorID,
		request.AssistantID,
		ReviewQueryOption{},
	)
	if err != nil {
		if err == exception.ErrNoRecord {
			return s.addOne(ctx, request)
		}

		return Review{}, err
	}

	now := mydate.Now()
	prev := existedReview.CreatedAt
	dayDiff := mydate.DaysDifference(now, prev)

	if dayDiff < 30 {
		return Review{}, exception.ErrBadRequest
	}

	return s.addOne(ctx, request)
}

func (s *reviewService) UpdateRecentOne(
	ctx inner.Context,
	request AddOneRequest,
) (
	Review,
	error,
) {
	_review, err := s.reviewRepository.FindOne_ByAuthorAndAssistantID(
		ctx,
		request.AuthorID,
		request.AssistantID,
		ReviewQueryOption{},
	)
	if err != nil {
		return Review{}, err
	}

	err = s.reviewRepository.UpdateOne_ByID(
		ctx,
		_review.ID,
		request.ToModelForInsert(),
	)
	if err != nil {
		return Review{}, err
	}

	return s.reviewRepository.FindOne_ByID(
		ctx,
		_review.ID,
		ReviewQueryOption{WithAuthor: true},
	)
}

func (s *reviewService) addOne(
	ctx inner.Context,
	request AddOneRequest,
) (
	Review,
	error,
) {
	_review, err := s.reviewRepository.InsertOne(
		ctx,
		request.ToModelForInsert(),
	)
	if err != nil {
		return Review{}, err
	}

	reviewWithAuthor, err := s.reviewRepository.FindOne_ByID(
		ctx,
		_review.ID,
		ReviewQueryOption{WithAuthor: true},
	)
	if err != nil {
		return Review{}, err
	}

	return reviewWithAuthor, nil
}

func NewReviewService(
	reviewRepository ReviewRepository,
) ReviewService {
	return &reviewService{
		reviewRepository: reviewRepository,
	}
}
