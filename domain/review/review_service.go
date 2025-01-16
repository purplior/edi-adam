package review

import (
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/lib/mydate"
)

type (
	ReviewService interface {
		WriteOne(
			ctx inner.Context,
			request WriteOneRequest,
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

func (s *reviewService) WriteOne(
	ctx inner.Context,
	request WriteOneRequest,
) (
	Review,
	error,
) {
	existedReview, err := s.reviewRepository.FindOne_ByAuthorAndAssistantID(
		ctx,
		request.AuthorID,
		request.AssistantID,
		ReviewJoinOption{},
	)
	if err != nil {
		if err == exception.ErrNoRecord {
			return s.reviewRepository.InsertOne(
				ctx,
				request.ToModelForInsert(),
			)
		} else {
			return Review{}, err
		}
	}
	if mydate.DaysDifference(mydate.Now(), existedReview.CreatedAt) < 31 {
		return Review{}, exception.ErrBadRequest
	}

	return s.reviewRepository.InsertOne(
		ctx,
		request.ToModelForInsert(),
	)
}

func NewReviewService(
	reviewRepository ReviewRepository,
) ReviewService {
	return &reviewService{
		reviewRepository: reviewRepository,
	}
}
