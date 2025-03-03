package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	domain "github.com/purplior/edi-adam/domain/missionlog"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	MissionLogController interface {
		Controller
		// 내 임무내역 가져오기
		GetMyList() common.Route
		// 내가 달성한 임무 보상받기
		ReceiveMine() common.Route
	}
)

type (
	missionLogController struct {
		missionLogService domain.MissionLogService
	}
)

func (c *missionLogController) GroupPath() string {
	return "/missionlogs"
}

func (c *missionLogController) Routes() []common.Route {
	return []common.Route{
		c.GetMyList(),
		c.ReceiveMine(),
	}
}

func (c *missionLogController) GetMyList() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/m/lst",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) error {
		mList, err := c.missionLogService.GetList(
			ctx.Session(),
			domain.QueryOption{
				UserID: ctx.Identity.ID,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			MissionLogs []model.MissionLog `json:"missionLogs"`
		}{
			MissionLogs: mList,
		})
	}

	return route
}

func (c *missionLogController) ReceiveMine() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/m/fst",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			ID uint `body:"id"`
		}
		if err = ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		if err = c.missionLogService.Receive(
			ctx.Session(),
			dto.ID,
			ctx.Identity.ID,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(nil)
	}

	return route
}

func NewMissionLogController(
	missionLogService domain.MissionLogService,
) MissionLogController {
	return &missionLogController{
		missionLogService: missionLogService,
	}
}
