package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	"github.com/podossaem/podoroot/domain/auth"
	domain "github.com/podossaem/podoroot/domain/me"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/domain/wallet"
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
		 *
		 */
		GetMyPodo() api.HandlerFunc
	}
)

type (
	meController struct {
		meService     domain.MeService
		authService   auth.AuthService
		userService   user.UserService
		walletService wallet.WalletService
		cm            inner.ContextManager
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

		userDetail, err := c.userService.GetDetailOneByID(
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

		wallet, err := c.walletService.GetOneByUserID(
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

func NewMeController(
	meService domain.MeService,
	authService auth.AuthService,
	userService user.UserService,
	walletService wallet.WalletService,
	cm inner.ContextManager,
) MeController {
	return &meController{
		meService:     meService,
		authService:   authService,
		userService:   userService,
		walletService: walletService,
		cm:            cm,
	}
}
