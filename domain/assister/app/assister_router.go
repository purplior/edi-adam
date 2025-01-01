package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
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

	assisterRouterGroup.GET(
		"/info-one",
		api.Handler(
			r.assisterController.GetInfoOne(),
			api.HandlerFuncOption{},
		),
	)

	assisterRouterGroup.POST(
		"/exec",
		api.Handler(
			r.assisterController.Execute(),
			api.HandlerFuncOption{},
		),
	)

	assisterRouterGroup.POST(
		"/exec-stream",
		api.Handler(
			r.assisterController.ExecuteAsStream(),
			api.HandlerFuncOption{},
		),
	)
}

func NewAssisterRouter(
	assisterController AssisterController,
) AssisterRouter {
	return &assisterRouter{
		assisterController: assisterController,
	}
}
