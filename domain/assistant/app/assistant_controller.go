package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/domain/user"
)

type (
	AssistantController interface {
		/**
		 * 신규 쌤비서 등록
		 */
		RegisterOne() api.HandlerFunc

		/**
		 * 포도쌤의 쌤비서 가져오기
		 */
		GetPodoList() api.HandlerFunc
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

func (c *assistantController) GetPodoList() api.HandlerFunc {
	return func(ctx *api.Context) error {
		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		assistants, err := c.assistantService.GetList(
			apiCtx,
			user.ID_Podo,
			true,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		assistantInfos := make([]domain.AssistantInfo, len(assistants))
		for i, assistant := range assistants {
			assistantInfos[i] = assistant.ToInfo()
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				AssistantInfos []domain.AssistantInfo `json:"assistantInfos"`
			}{
				AssistantInfos: assistantInfos,
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
