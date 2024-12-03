package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
)

type (
	AssisterFormRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	assisterFormRouter struct {
		assisterFormController AssisterFormController
	}
)

func (r *assisterFormRouter) Attach(router *echo.Group) {
	assisterFormRouterGroup := router.Group("/assisterforms")

	assisterFormRouterGroup.POST(
		"/",
		api.Handler(
			r.assisterFormController.RegisterOne(),
			api.HandlerFuncOption{},
		),
	)

	assisterFormRouterGroup.GET(
		"/one",
		api.Handler(
			r.assisterFormController.GetViewOne(),
			api.HandlerFuncOption{},
		),
	)

	assisterFormRouterGroup.GET(
		"/admin/one",
		api.Handler(
			r.assisterFormController.GetOne_ForAdmin(),
			api.HandlerFuncOption{AdminOnly: true},
		),
	)

	assisterFormRouterGroup.PUT(
		"/admin/one",
		api.Handler(
			r.assisterFormController.PutOne_ForAdmin(),
			api.HandlerFuncOption{AdminOnly: true},
		),
	)

	assisterFormRouterGroup.POST(
		"/admin/one",
		api.Handler(
			r.assisterFormController.CreateOne_ForAdmin(),
			api.HandlerFuncOption{AdminOnly: true},
		),
	)
}

func NewAssisterFormRouter(
	assisterFormController AssisterFormController,
) AssisterFormRouter {
	return &assisterFormRouter{
		assisterFormController: assisterFormController,
	}
}
