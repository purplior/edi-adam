package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
)

type (
	MeRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	meRouter struct {
		meController MeController
	}
)

func (r *meRouter) Attach(router *echo.Group) {
	meRouterGroup := router.Group("/me")

	meRouterGroup.GET(
		"/identity",
		api.Handler(
			r.meController.GetMyIdentity(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.GET(
		"/detail",
		api.Handler(
			r.meController.GetMyDetail(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.GET(
		"/podo",
		api.Handler(
			r.meController.GetMyPodo(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.GET(
		"/assistant-infos",
		api.Handler(
			r.meController.GetMyAssistantInfos(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.GET(
		"/temp/at",
		api.Handler(
			r.meController.GetTempAccessToken(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.GET(
		"/assistant/by_view/:view_id",
		api.Handler(
			r.meController.GetMyAssistant(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.GET(
		"/assistant/detail_by_view/:view_id",
		api.Handler(
			r.meController.GetMyAssistantDetail(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.POST(
		"/assistant",
		api.Handler(
			r.meController.RegisterMyAssistant(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.PATCH(
		"/assistant/:id",
		api.Handler(
			r.meController.UpdateMyAssistant(),
			api.HandlerFuncOption{},
		),
	)

	meRouterGroup.DELETE(
		"/assistant/:id",
		api.Handler(
			r.meController.RemoveMyAssistant(),
			api.HandlerFuncOption{},
		),
	)
}

func NewMeRouter(
	meController MeController,
) MeRouter {
	return &meRouter{
		meController: meController,
	}
}
