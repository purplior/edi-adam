package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
)

type (
	MissionRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	missionRouter struct {
		missionController MissionController
	}
)

func (r *missionRouter) Attach(router *echo.Group) {
	missionRouterGroup := router.Group("/missions")

	missionRouterGroup.GET("/plist", api.Handler(
		r.missionController.GetPaginatedList(),
		api.HandlerFuncOption{},
	))

	missionRouterGroup.POST("/receive", api.Handler(
		r.missionController.ReceiveOne(),
		api.HandlerFuncOption{},
	))
}

func NewMissionRouter(
	missionController MissionController,
) MissionRouter {
	return &missionRouter{
		missionController: missionController,
	}
}
