package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/mission"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	MissionController interface {
		GetPaginatedList() api.HandlerFunc

		ReceiveOne() api.HandlerFunc
	}
)

type (
	missionController struct {
		missionService domain.MissionService
		cm             inner.ContextManager
	}
)

func (c *missionController) GetPaginatedList() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrBadRequest)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		missions, meta, err := c.missionService.GetPaginatedList_ByUserID(
			innerCtx,
			ctx.Identity.ID,
			dt.Int(ctx.QueryParam("p")),
			dt.Int(ctx.QueryParam("ps")),
		)

		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Missions []domain.Mission          `json:"missions"`
				Meta     pagination.PaginationMeta `json:"meta"`
			}{
				Missions: missions,
				Meta:     meta,
			},
		})
	}
}

func (c *missionController) ReceiveOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrBadRequest)
		}

		var dto struct {
			ID string `json:"id"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		err := c.missionService.ReceiveOne(
			innerCtx,
			dto.ID,
			ctx.Identity.ID,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: nil,
		})
	}
}

func NewMissionController(
	missionService domain.MissionService,
	cm inner.ContextManager,
) MissionController {
	return &missionController{
		missionService: missionService,
		cm:             cm,
	}
}
