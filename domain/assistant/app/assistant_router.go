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
			api.HandlerFuncOption{},
		),
	)

	assistantRouterGroup.GET(
		"/detail/:assistant_view_id",
		api.Handler(
			r.assistantController.GetDetailOne(),
			api.HandlerFuncOption{},
		),
	)

	assistantRouterGroup.GET(
		"/podo-list",
		api.Handler(
			r.assistantController.GetPodoInfoList(),
			api.HandlerFuncOption{},
		),
	)

	assistantRouterGroup.GET(
		"/admin/pages",
		api.Handler(
			r.assistantController.GetPaginatedList_ForAdmin(),
			api.HandlerFuncOption{AdminOnly: true},
		),
	)

	assistantRouterGroup.GET(
		"/admin/one",
		api.Handler(
			r.assistantController.GetOne_ForAdmin(),
			api.HandlerFuncOption{AdminOnly: true},
		),
	)

	assistantRouterGroup.PUT(
		"/admin/one",
		api.Handler(
			r.assistantController.PutOne_ForAdmin(),
			api.HandlerFuncOption{AdminOnly: true},
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
