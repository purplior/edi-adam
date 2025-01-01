package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
)

type (
	AssistantRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	assistantRouter struct {
		assistantController AssistantController
	}
)

func (r *assistantRouter) Attach(router *echo.Group) {
	assistantRouterGroup := router.Group("/assistants")

	assistantRouterGroup.GET(
		"/detail/one",
		api.Handler(
			r.assistantController.GetDetailOne(),
			api.HandlerFuncOption{},
		),
	)

	assistantRouterGroup.GET(
		"/info-list",
		api.Handler(
			r.assistantController.GetInfoList(),
			api.HandlerFuncOption{},
		),
	)
}

func NewAssistantRouter(
	assistantController AssistantController,
) AssistantRouter {
	return &assistantRouter{
		assistantController: assistantController,
	}
}
