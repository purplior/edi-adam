package app

import (
	"github.com/purplior/sbec/application/api"
	"github.com/purplior/sbec/application/response"
	domain "github.com/purplior/sbec/domain/mission"
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/shared/pagination"
	"github.com/purplior/sbec/lib/dt"
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

		missions, meta, err := c.missionService.GetPaginatedList_OnlyPublic_ByUserID(
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
