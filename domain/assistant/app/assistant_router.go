package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
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

	assistantRouterGroup.POST(
		"/",
		api.Handler(
			r.assistantController.RegisterOne(),
		),
	)

	assistantRouterGroup.GET(
		"/:assistant_id",
		api.Handler(
			r.assistantController.GetDetailOne(),
		),
	)

	assistantRouterGroup.GET(
		"/list/podo",
		api.Handler(
			r.assistantController.GetPodoInfoList(),
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
