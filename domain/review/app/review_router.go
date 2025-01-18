package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
)

type (
	ReviewRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	reviewRouter struct {
		reviewController ReviewController
	}
)

func (r *reviewRouter) Attach(router *echo.Group) {
	rg := router.Group("/reviews")

	rg.GET("/info-plist", api.Handler(
		r.reviewController.GetInfoPaginatedList(),
		api.HandlerFuncOption{},
	))
}

func NewReviewRouter(
	reviewController ReviewController,
) ReviewRouter {
	return &reviewRouter{
		reviewController: reviewController,
	}
}
