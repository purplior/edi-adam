package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/review"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
)

type (
	ReviewController interface {
		WriteOne() api.HandlerFunc
	}
)

type (
	reviewController struct {
		reviewService domain.ReviewService
		cm            inner.ContextManager
	}
)

func (c *reviewController) WriteOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}
		var request domain.WriteOneRequest
		if err := ctx.Bind(&request); err != nil {
			return ctx.SendError(err)
		}

		if request.Score <= 0 || request.Score > 5 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		request.AuthorID = ctx.Identity.ID
		review, err := c.reviewService.WriteOne(
			innerCtx,
			request,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Review domain.Review `json:"review"`
			}{
				Review: review,
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
