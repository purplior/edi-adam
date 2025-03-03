package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	domain "github.com/purplior/edi-adam/domain/mission"
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	MissionController interface {
		Controller
		// 미션 목록 가져오기
		GetPaginatedList() common.Route
	}
)

type (
	missionController struct {
		missionService domain.MissionService
	}
)

func (c *missionController) GroupPath() string {
	return "/missions"
}

func (c *missionController) Routes() []common.Route {
	return []common.Route{
		c.GetPaginatedList(),
	}
}

func (c *missionController) GetPaginatedList() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/o/lst",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			pagination.PaginationRequest
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		mList, meta, err := c.missionService.GetPaginatedList(
			ctx.Session(),
			pagination.PaginationQuery[domain.QueryOption]{
				QueryOption: domain.QueryOption{
					IsPublic: true,
				},
				PageRequest: dto.PaginationRequest,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Missions []model.Mission           `json:"missions"`
			Meta     pagination.PaginationMeta `json:"meta"`
		}{
			Missions: mList,
			Meta:     meta,
		})
	}

	return route
}

func NewMissionController(
	missionService domain.MissionService,
) MissionController {
	return &missionController{
		missionService: missionService,
	}
}
