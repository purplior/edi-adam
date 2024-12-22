package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	"github.com/purplior/podoroot/domain/assistant"
	"github.com/purplior/podoroot/domain/auth"
	domain "github.com/purplior/podoroot/domain/me"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/user"
	"github.com/purplior/podoroot/domain/wallet"
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
		 * 내 어시 등록하기
		 */
		RegisterMyAssistant() api.HandlerFunc

		/**
		 * 내 어시 제거하기
		 */
		RemoveMyAssistant() api.HandlerFunc
	}
)

type (
	meController struct {
		meService        domain.MeService
		assistantService assistant.AssistantService
		authService      auth.AuthService
		userService      user.UserService
		walletService    wallet.WalletService
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

func (c *meController) RegisterMyAssistant() api.HandlerFunc {
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

func (c *meController) RemoveMyAssistant() api.HandlerFunc {
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

func NewMeController(
	meService domain.MeService,
	assistantService assistant.AssistantService,
	authService auth.AuthService,
	userService user.UserService,
	walletService wallet.WalletService,
	cm inner.ContextManager,
) MeController {
	return &meController{
		meService:        meService,
		assistantService: assistantService,
		authService:      authService,
		userService:      userService,
		walletService:    walletService,
		cm:               cm,
	}
}
