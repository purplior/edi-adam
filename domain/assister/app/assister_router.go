package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
)

type (
	AssisterRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	assisterRouter struct {
		assisterController AssisterController
	}
)

func (r *assisterRouter) Attach(router *echo.Group) {
	assisterRouterGroup := router.Group("/assisters")

	assisterRouterGroup.POST("/exec", api.Handler(
		r.assisterController.Execute(),
		api.HandlerFuncOption{},
	))

	assisterRouterGroup.GET("/admin/one", api.Handler(
		r.assisterController.GetOne_ForAdmin(),
		api.HandlerFuncOption{AdminOnly: true},
	))

	assisterRouterGroup.GET("/admin/pages", api.Handler(
		r.assisterController.GetPaginatedList_ForAdmin(),
		api.HandlerFuncOption{AdminOnly: true},
	))
}

func NewAssisterRouter(
	assisterController AssisterController,
) AssisterRouter {
	return &assisterRouter{
		assisterController: assisterController,
	}
}
