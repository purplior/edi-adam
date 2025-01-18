package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/review"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
)

type (
	ReviewController interface {
		GetInfoPaginatedList() api.HandlerFunc
	}
)

type (
	reviewController struct {
		reviewService domain.ReviewService
		cm            inner.ContextManager
	}
)

func (c *reviewController) GetInfoPaginatedList() api.HandlerFunc {
	return func(ctx *api.Context) error {
		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assistantID := ctx.QueryParam("assistant_id")
		if len(assistantID) < 1 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		pageRequest, err := ctx.PaginationRequest()
		if err != nil {
			return ctx.SendError(err)
		}

		reviews, pageMeta, err := c.reviewService.GetPaginatedList_ByAssistantID(
			innerCtx,
			assistantID,
			pageRequest,
		)
		if err != nil && err != exception.ErrNoRecord {
			return ctx.SendError(err)
		}

		reviewInfos := make([]domain.ReviewInfo, len(reviews))
		for i, review := range reviews {
			reviewInfos[i] = review.ToInfo()
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				ReviewInfos []domain.ReviewInfo       `json:"reviewInfos"`
				PageMeta    pagination.PaginationMeta `json:"pageMeta"`
			}{
				ReviewInfos: reviewInfos,
				PageMeta:    pageMeta,
			},
		})
	}
}

func NewReviewController(
	reviewService domain.ReviewService,
	cm inner.ContextManager,
) ReviewController {
	return &reviewController{
		reviewService: reviewService,
		cm:            cm,
	}
}
