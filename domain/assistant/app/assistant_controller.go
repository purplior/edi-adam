package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/shared/pagination"
	"github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/lib/dt"
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

		/**
		 * 카테고리별 쌤비서 가져오기
		 */
		GetPaginatedList_ForAdmin() api.HandlerFunc
	}
)

type (
	assistantController struct {
		assistantService domain.AssistantService
		cm               inner.ContextManager
	}
)

func (c *assistantController) RegisterOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.RegisterOneRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assistant, err := c.assistantService.RegisterOne(
			innerCtx,
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
		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assistantViewID := ctx.Param("assistant_view_id")
		if len(assistantViewID) == 0 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		assistantDetail, err := c.assistantService.GetDetailOne_ByViewID(
			innerCtx,
			assistantViewID,
			domain.AssistantJoinOption{
				WithAuthor:    true,
				WithAssisters: true,
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
		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assistantInfos, err := c.assistantService.GetInfoList_ByAuthor(
			innerCtx,
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

func (c *assistantController) GetPaginatedList_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		authorID := ctx.QueryParam("author_id")
		page := dt.Int(ctx.QueryParam("p"))
		pageSize := dt.Int(ctx.QueryParam("ps"))

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assistants, meta, err := c.assistantService.GetPaginatedList_ByAuthor(
			innerCtx,
			authorID,
			page,
			pageSize,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Assistants []domain.Assistant        `json:"assistants"`
				Meta       pagination.PaginationMeta `json:"meta"`
			}{
				Assistants: assistants,
				Meta:       meta,
			},
		})
	}
}

func NewAssistantController(
	assistantService domain.AssistantService,
	cm inner.ContextManager,
) AssistantController {
	return &assistantController{
		assistantService: assistantService,
		cm:               cm,
	}
}
