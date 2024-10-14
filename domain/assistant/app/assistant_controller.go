package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/context"
)

type (
	AssistantController interface {
		/**
		 * 신규 어시스턴트 등록
		 */
		RegisterOne() api.HandlerFunc
	}
)

type (
	assistantController struct {
		assistantService domain.AssistantService
	}
)

func (c *assistantController) RegisterOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.RegisterOneRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		assistant, err := c.assistantService.RegisterOne(
			apiCtx,
			ctx.Identity.ID,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Assistant domain.Assistant `json:"assistant"`
			}{
				Assistant: assistant,
			},
		})
	}
}

func NewAssistantController(
	assistantService domain.AssistantService,
) AssistantController {
	return &assistantController{
		assistantService: assistantService,
	}
}
