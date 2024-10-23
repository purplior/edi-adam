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

	assisterRouterGroup.POST("/:assister_id/exec", api.Handler(
		r.assisterController.Execute(),
	))
}

func NewAssisterRouter(
	assisterController AssisterController,
) AssisterRouter {
	return &assisterRouter{
		assisterController: assisterController,
	}
}
