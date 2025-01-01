package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
)

type (
	CategoryRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	categoryRouter struct {
		categoryController CategoryController
	}
)

func (r *categoryRouter) Attach(router *echo.Group) {
	categoryRouterGroup := router.Group("/categories")

	categoryRouterGroup.GET(
		"/main",
		api.Handler(
			r.categoryController.GetMainCategoryInfoList(),
			api.HandlerFuncOption{},
		),
	)
}

func NewCategoryRouter(
	categoryController CategoryController,
) CategoryRouter {
	return &categoryRouter{
		categoryController: categoryController,
	}
}
