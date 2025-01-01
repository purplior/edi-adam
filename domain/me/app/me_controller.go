package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	"github.com/purplior/podoroot/domain/assistant"
	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/auth"
	"github.com/purplior/podoroot/domain/bookmark"
	domain "github.com/purplior/podoroot/domain/me"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/domain/user"
	"github.com/purplior/podoroot/domain/wallet"
	"github.com/purplior/podoroot/lib/dt"
	"github.com/purplior/podoroot/lib/validator"
)

type (
	MeController interface {
		/**
		 * 내 정보(인증에 포함된 간단한 식별정보) 가져오기
		 */
		GetMyIdentity() api.HandlerFunc

		/**
		 * 내 정보 가져오기
		 */
		GetMyDetail() api.HandlerFunc

		/**
		 * 나의 임시 액세스 토큰 발급하기 (유효기간 1시간)
		 */
		GetTempAccessToken() api.HandlerFunc

		/**
		 * 나의 포도 확인하기
		 */
		GetMyPodo() api.HandlerFunc

		/**
		 * 나의 어시목록 확인하기
		 */
		GetMyAssistantInfoPaginatedList() api.HandlerFunc

		/**
		 * 나의 어시 확인하기
		 */
		GetMyAssistantOne() api.HandlerFunc

		/**
		 * 나의 어시 상세 확인하기
		 */
		GetMyAssistantDetailOne() api.HandlerFunc

		/**
		 * 나의 북마크 확인하기
		 */
		GetMyBookmarkOne() api.HandlerFunc

		/**
		 * 나의 북마크 목록 확인하기
		 */
		GetMyBookmarkPaginatedList() api.HandlerFunc

		/**
		 * 내 어시 등록하기
		 */
		RegisterMyAssistantOne() api.HandlerFunc

		/**
		 * 내 어시 수정하기
		 */
		UpdateMyAssistantOne() api.HandlerFunc

		/**
		 * 내 어시 제거하기
		 */
		RemoveMyAssistantOne() api.HandlerFunc

		/**
		 * 북마크 토글하기
		 */
		ToggleBookmarkOne() api.HandlerFunc
	}
)

type (
	meController struct {
		meService        domain.MeService
		assistantService assistant.AssistantService
		assisterService  assister.AssisterService
		authService      auth.AuthService
		userService      user.UserService
		walletService    wallet.WalletService
		bookmarkService  bookmark.BookmarkService
		cm               inner.ContextManager
	}
)

func (c *meController) GetMyIdentity() api.HandlerFunc {
	return func(ctx *api.Context) error {
		identity := ctx.Identity

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Ok,
			Data:   identity,
		})
	}
}

func (c *meController) GetMyDetail() api.HandlerFunc {
	return func(ctx *api.Context) error {
		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		userDetail, err := c.userService.GetDetailOne_ByID(
			innerCtx,
			ctx.Identity.ID,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				UserDetail user.UserDetail `json:"userDetail"`
			}{
				UserDetail: userDetail,
			},
		})
	}
}

func (c *meController) GetTempAccessToken() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		accessToken, err := c.authService.GetTempAccessToken(innerCtx, *ctx.Identity)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				TempAccessToken string `json:"tempAccessToken"`
			}{
				TempAccessToken: accessToken,
			},
		})
	}
}

func (c *meController) GetMyPodo() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		wallet, err := c.walletService.GetOne_ByUserID(
			innerCtx,
			ctx.Identity.ID,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Podo int `json:"podo"`
			}{
				Podo: wallet.Podo,
			},
		})
	}
}

func (c *meController) GetMyAssistantInfoPaginatedList() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		page := dt.Int(ctx.QueryParam("p"))
		if page < 1 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assistants, pageMeta, err := c.assistantService.GetPaginatedList_ByAuthor(
			innerCtx,
			ctx.Identity.ID,
			pagination.PaginationRequest{
				Page: page,
				Size: 16,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		assistantInfos := make([]assistant.AssistantInfo, len(assistants))
		for i, assistant := range assistants {
			assistantInfos[i] = assistant.ToInfo()
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				AssistantInfos []assistant.AssistantInfo `json:"assistantInfos"`
				PageMeta       pagination.PaginationMeta `json:"pageMeta"`
			}{
				AssistantInfos: assistantInfos,
				PageMeta:       pageMeta,
			},
		})
	}
}

func (c *meController) GetMyAssistantOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		var _assistant assistant.Assistant
		var err error = nil

		id := ctx.QueryParam("id")
		viewID := ctx.QueryParam("view_id")

		if len(id) > 0 {
			_assistant, err = c.assistantService.GetOne_ByID(
				innerCtx,
				id,
				assistant.AssistantJoinOption{
					WithCategory: true,
					WithAssister: true,
				},
			)
		} else if len(viewID) > 0 {
			_assistant, err = c.assistantService.GetOne_ByViewID(
				innerCtx,
				viewID,
				assistant.AssistantJoinOption{
					WithCategory: true,
					WithAssister: true,
				},
			)
		}
		if err != nil {
			return ctx.SendError(err)
		}
		if _assistant.AuthorID != ctx.Identity.ID {
			return ctx.SendError(exception.ErrBadRequest)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Assistant assistant.Assistant `json:"assistant"`
			}{
				Assistant: _assistant,
			},
		})
	}
}

func (c *meController) GetMyAssistantDetailOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		var _assistant assistant.Assistant
		var err error = nil

		id := ctx.QueryParam("id")
		viewID := ctx.QueryParam("view_id")

		if len(id) > 0 {
			_assistant, err = c.assistantService.GetOne_ByID(
				innerCtx,
				id,
				assistant.AssistantJoinOption{
					WithCategory: true,
					WithAssister: true,
				},
			)
		} else if len(viewID) > 0 {
			_assistant, err = c.assistantService.GetOne_ByViewID(
				innerCtx,
				viewID,
				assistant.AssistantJoinOption{
					WithCategory: true,
					WithAssister: true,
				},
			)
		}
		if err != nil {
			return ctx.SendError(err)
		}
		if _assistant.AuthorID != ctx.Identity.ID {
			return ctx.SendError(exception.ErrBadRequest)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				AssistantDetail assistant.AssistantDetail `json:"assistantDetail"`
			}{
				AssistantDetail: _assistant.ToDetail(),
			},
		})
	}
}

func (c *meController) GetMyBookmarkOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assistantID := ctx.QueryParam("assistant_id")
		if len(assistantID) == 0 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		bookmark, err := c.bookmarkService.GetOne_ByUserIDAndAssistantID(
			innerCtx,
			ctx.Identity.ID,
			assistantID,
		)
		if err != nil && err != exception.ErrNoRecord {
			return ctx.SendError(err)
		}

		isEmpty := len(bookmark.ID) == 0 || err == exception.ErrNoRecord

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Exist bool `json:"exist"`
			}{
				Exist: !isEmpty,
			},
		})
	}
}

func (c *meController) GetMyBookmarkPaginatedList() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		page := dt.Int(ctx.QueryParam("p"))
		pageSize := dt.Int(ctx.QueryParam("ps"))

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		bookmarks, pageMeta, err := c.bookmarkService.GetPaginatedList_ByUserID(
			innerCtx,
			ctx.Identity.ID,
			pagination.PaginationRequest{
				Page: page,
				Size: pageSize,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Bookmarks []bookmark.Bookmark       `json:"bookmarks"`
				PageMeta  pagination.PaginationMeta `json:"pageMeta"`
			}{
				Bookmarks: bookmarks,
				PageMeta:  pageMeta,
			},
		})
	}
}

func (c *meController) RegisterMyAssistantOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		var request assistant.RegisterOneRequest
		if err := ctx.Bind(&request); err != nil {
			return ctx.SendError(err)
		}

		if err := validator.CheckValidAssistantRegisterRequest(request); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		if _, err := c.assistantService.RegisterOne(
			innerCtx,
			ctx.Identity.ID,
			request,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Created,
		})
	}
}

func (c *meController) UpdateMyAssistantOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		assistantID := ctx.Param("id")
		if len(assistantID) == 0 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		var request assistant.UpdateOneRequest
		if err := ctx.Bind(&request); err != nil {
			return ctx.SendError(err)
		}

		request.ID = assistantID

		if err := validator.CheckValidAssistantUpdateRequest(request); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		if err := c.assistantService.UpdateOne(
			innerCtx,
			ctx.Identity.ID,
			request,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func (c *meController) RemoveMyAssistantOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		assistantID := ctx.Param("id")
		if len(assistantID) == 0 {
			return ctx.SendError(exception.ErrBadRequest)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		if err := c.assistantService.RemoveOne_ByID(
			innerCtx,
			ctx.Identity.ID,
			assistantID,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func (c *meController) ToggleBookmarkOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		var dto struct {
			AssistantID string `json:"assistantId"`
			Toggle      bool   `json:"toggle"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		userID := ctx.Identity.ID

		if err := c.bookmarkService.ToggleOne(
			innerCtx,
			userID,
			dto.AssistantID,
			dto.Toggle,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func NewMeController(
	meService domain.MeService,
	assistantService assistant.AssistantService,
	assisterService assister.AssisterService,
	authService auth.AuthService,
	userService user.UserService,
	walletService wallet.WalletService,
	bookmarkService bookmark.BookmarkService,
	cm inner.ContextManager,
) MeController {
	return &meController{
		meService:        meService,
		assistantService: assistantService,
		assisterService:  assisterService,
		authService:      authService,
		userService:      userService,
		walletService:    walletService,
		bookmarkService:  bookmarkService,
		cm:               cm,
	}
}
