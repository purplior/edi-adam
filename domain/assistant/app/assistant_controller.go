package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/assistant"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	AssistantController interface {
		/**
		 * 쌤비서 상세정보 가져오기
		 */
		GetDetailOne() api.HandlerFunc

		/**
		 * 쌤비서 목록 가져오기 (카테고리)
		 */
		GetInfoList_ByCategory() api.HandlerFunc

		/**
		 * 카테고리별 쌤비서 가져오기
		 */
		GetPaginatedList_ForAdmin() api.HandlerFunc

		/**
		 * 쌤비서 가져오기 (어드민용)
		 */
		GetOne_ForAdmin() api.HandlerFunc

		/**
		 * 쌤비서 수정하기 (어드민용)
		 */
		PutOne_ForAdmin() api.HandlerFunc

		/**
		 * 쌤비서 수정하기 (어드민용)
		 */
		CreateOne_ForAdmin() api.HandlerFunc
	}
)

type (
	assistantController struct {
		assistantService domain.AssistantService
		cm               inner.ContextManager
	}
)

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

func (c *assistantController) GetInfoList_ByCategory() api.HandlerFunc {
	return func(ctx *api.Context) error {
		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		categoryAlias := ctx.QueryParam("c")

		assistantInfos, err := c.assistantService.GetInfoList_ByCategory(
			innerCtx,
			categoryAlias,
			domain.AssistantJoinOption{
				WithAuthor:   true,
				WithCategory: true,
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

func (c *assistantController) GetOne_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		id := ctx.QueryParam("id")

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assistant, err := c.assistantService.GetOne_ByID(
			innerCtx,
			id,
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
				Assistant domain.Assistant `json:"assistant"`
			}{
				Assistant: assistant,
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

func (c *assistantController) PutOne_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			Assistant domain.Assistant `json:"assistant"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		if err := c.assistantService.PutOne(
			innerCtx,
			dto.Assistant,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func (c *assistantController) CreateOne_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			Assistant domain.Assistant `json:"assistant"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		if err := c.assistantService.CreateOne(
			innerCtx,
			dto.Assistant,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Created,
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
