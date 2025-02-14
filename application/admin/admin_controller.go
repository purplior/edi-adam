package admin

import (
	"github.com/purplior/sbec/application/api"
	"github.com/purplior/sbec/application/response"
	"github.com/purplior/sbec/domain/assistant"
	"github.com/purplior/sbec/domain/shared/inner"
)

type (
	AdminController interface {
		ApproveAssistantOne() api.HandlerFunc
	}
)

type (
	adminController struct {
		assistantService assistant.AssistantService
		cm               inner.ContextManager
	}
)

func (c *adminController) ApproveAssistantOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		var dto struct {
			ID       string   `json:"id"`
			Cost     int      `json:"cost"`
			MetaTags []string `json:"metaTags"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		if err := c.assistantService.ApproveOne(
			innerCtx,
			dto.ID,
			dto.Cost,
			dto.MetaTags,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func NewAdminController(
	assistantService assistant.AssistantService,
	cm inner.ContextManager,
) AdminController {
	return &adminController{
		assistantService: assistantService,
		cm:               cm,
	}
}
