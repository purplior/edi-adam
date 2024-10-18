package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/user"
)

type (
	AssistantController interface {
		/**
		 * 신규 쌤비서 등록
		 */
		RegisterOne() api.HandlerFunc

		/**
		 * 쌤비서 상세정보 가져오기
		 */
		GetDetailOne() api.HandlerFunc

		/**
		 * 포도쌤의 쌤비서 가져오기
		 */
		GetPodoInfoList() api.HandlerFunc
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

func (c *assistantController) GetDetailOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		assistantViewID := ctx.Param("assistant_view_id")
		if len(assistantViewID) == 0 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		assistantDetail, err := c.assistantService.GetDetailOneByViewID(
			apiCtx,
			assistantViewID,
			domain.AssistantJoinOption{
				WithAuthor:          true,
				WithDefaultAssister: true,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				AssistantDetail domain.AssistantDetail `json:"assistantDetail"`
			}{
				AssistantDetail: assistantDetail,
			},
		})
	}
}

func (c *assistantController) GetPodoInfoList() api.HandlerFunc {
	return func(ctx *api.Context) error {
		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		assistantInfos, err := c.assistantService.GetInfoListByAuthor(
			apiCtx,
			user.ID_Podo,
			domain.AssistantJoinOption{
				WithAuthor: true,
			},
		)
		if err != nil {
			return ctx.SendError(err)
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
