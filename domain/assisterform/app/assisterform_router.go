package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
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
			r.assisterFormController.GetOne_ByAssisterID_ForAdmin(),
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
