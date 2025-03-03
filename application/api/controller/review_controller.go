package controller

import (
	"net/http"
	"unicode/utf8"

	"github.com/purplior/edi-adam/application/common"
	"github.com/purplior/edi-adam/application/response"
	domain "github.com/purplior/edi-adam/domain/review"
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	ReviewController interface {
		Controller
		// 리뷰 카드 가져오기
		GetPaginatedList_Card() common.Route
		// 리뷰 남기기
		WriteMine() common.Route
	}
)

type (
	reviewController struct {
		reviewService domain.ReviewService
	}
)

func (c *reviewController) GroupPath() string {
	return "/reviews"
}

func (c *reviewController) Routes() []common.Route {
	return []common.Route{
		c.GetPaginatedList_Card(),
		c.WriteMine(),
	}
}

func (c *reviewController) GetPaginatedList_Card() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/o/lst",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			AssistantID uint `query:"assistant_id,required"`
			pagination.PaginationRequest
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		reviews, pageMeta, err := c.reviewService.GetPaginatedList(
			ctx.Session(),
			pagination.PaginationQuery[domain.QueryOption]{
				QueryOption: domain.QueryOption{
					AssistantID: dto.AssistantID,
				},
				PageRequest: dto.PaginationRequest,
			},
		)
		if err != nil && err != exception.ErrNoRecord {
			return ctx.SendError(err)
		}

		reviewCards := make([]model.ReviewCard, len(reviews))
		for i, review := range reviews {
			reviewCards[i] = review.ToCard()
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				ReviewCards []model.ReviewCard        `json:"reviewCards"`
				PageMeta    pagination.PaginationMeta `json:"pageMeta"`
			}{
				ReviewCards: reviewCards,
				PageMeta:    pageMeta,
			},
		})
	}

	return route
}

func (c *reviewController) WriteMine() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/m/lst",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			AssistantID uint    `body:"assistantId,required"`
			Content     string  `body:"content,required"`
			Score       float64 `body:"score,required"`
		}
		if err = ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}
		if dto.Score <= 0 || dto.Score > 5 {
			return ctx.SendCustomError(
				response.Status_BadRequest,
				"별점을 다시 확인해주세요.",
			)
		}
		contentLen := utf8.RuneCountInString(dto.Content)
		if contentLen < 20 || contentLen > 500 {
			return ctx.SendCustomError(
				response.Status_BadRequest,
				"길이는 20자 이상 500자 이하로 작성해주세요.",
			)
		}

		m, err := c.reviewService.Add(
			ctx.Session(),
			domain.AddDTO{
				AssistantID: dto.AssistantID,
				Content:     dto.Content,
				Score:       dto.Score,
			},
		)
		if err != nil {
			if err == exception.ErrBadRequest {
				return ctx.SendCustomError(
					response.Status_Unprocessable,
					"동일한 어시에는 한달에 한번만 리뷰를 남길 수 있어요.",
				)
			}
			return ctx.SendError(err)
		}

		return ctx.SendCreated(struct {
			ReviewCard model.ReviewCard `json:"reviewCard"`
		}{
			ReviewCard: m.ToCard(),
		})
	}

	return route
}

func NewReviewController(
	reviewService domain.ReviewService,
) ReviewController {
	return &reviewController{
		reviewService: reviewService,
	}
}
