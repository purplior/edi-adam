package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	"github.com/podossaem/podoroot/domain/auth"
	domain "github.com/podossaem/podoroot/domain/me"
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/domain/shared/exception"
)

type (
	MeController interface {
		/**
		 * 내 정보 가져오기
		 */
		GetMyIdentity() api.HandlerFunc

		/**
		 * 나의 임시 액세스 토큰 발급하기 (유효기간 1시간)
		 */
		GetTempAccessToken() api.HandlerFunc
	}
)

type (
	meController struct {
		meService   domain.MeService
		authService auth.AuthService
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

func (c *meController) GetTempAccessToken() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		accessToken, err := c.authService.GetTempAccessToken(apiCtx, *ctx.Identity)
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

func NewMeController(
	meService domain.MeService,
	authService auth.AuthService,
) MeController {
	return &meController{
		meService:   meService,
		authService: authService,
	}
}
