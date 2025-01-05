package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/assistant"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
)

type (
	AssistantController interface {
		/**
		 * 어시 상세정보 가져오기
		 */
		GetDetailOne() api.HandlerFunc

		/**
		 * 어시 목록 가져오기 (카테고리)
		 */
		GetInfoPaginatedList() api.HandlerFunc
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

		viewID := ctx.QueryParam("view_id")
		if len(viewID) == 0 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		_assistant, err := c.assistantService.GetOne_ByViewID(
			innerCtx,
			viewID,
			domain.AssistantJoinOption{
				WithAuthor:   true,
				WithAssister: true,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}
		if !_assistant.IsPublic {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				AssistantDetail domain.AssistantDetail `json:"assistantDetail"`
			}{
				AssistantDetail: _assistant.ToDetail(),
			},
		})
	}
}

func (c *assistantController) GetInfoPaginatedList() api.HandlerFunc {
	return func(ctx *api.Context) error {
		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		categoryAlias := ctx.QueryParam("c")
		if len(categoryAlias) == 0 {
			return ctx.SendError(exception.ErrBadRequest)
		}
		pageRequest := ctx.PaginationRequest()
		assistants, pageMeta, err := c.assistantService.GetPaginatedList_ByCategoryAlias(
			innerCtx,
			categoryAlias,
			true,
			pageRequest,
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
				AssistantInfos []domain.AssistantInfo    `json:"assistantInfos"`
				PageMeta       pagination.PaginationMeta `json:"pageMeta"`
			}{
				AssistantInfos: assistantInfos,
				PageMeta:       pageMeta,
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
