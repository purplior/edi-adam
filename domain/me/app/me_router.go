package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
)

type (
	MeRouter interface {
		Attach(router *echo.Group)
	}

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
		),
	)

	meRouterGroup.GET(
		"/temp/at",
		api.Handler(
			r.meController.GetTempAccessToken(),
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
