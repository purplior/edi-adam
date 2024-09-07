package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/auth"
	"github.com/podossaem/podoroot/domain/context"
)

type (
	AuthController interface {
		/**
		 * 이메일로 로그인
		 */
		SignInByEmailVerification() api.HandlerFunc

		/**
		 * 이메일로 회원가입
		 */
		SignUpByEmailVerification() api.HandlerFunc

		/**
		 * 토큰 재발급
		 */
		RefreshIdentityToken() api.HandlerFunc
	}
)

type (
	authController struct {
		authService domain.AuthService
	}
)

func (c *authController) SignInByEmailVerification() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.SignInByEmailVerificationRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		token, err := c.authService.SignInByEmailVerification(
			apiCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Token domain.IdentityToken `json:"token"`
			}{
				Token: token,
			},
		})
	}
}

func (c *authController) SignUpByEmailVerification() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.SignUpByEmailVerificationRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		token, err := c.authService.SignUpByEmailVerification(
			apiCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Token domain.IdentityToken `json:"token"`
			}{
				Token: token,
			},
		})
	}
}

func (c *authController) RefreshIdentityToken() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.IdentityToken
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		newIdentityToken, err := c.authService.RefreshIdentityToken(apiCtx, dto)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Ok,
			Data:   newIdentityToken,
		})
	}
}

func NewAuthController(
	authService domain.AuthService,
) AuthController {
	return &authController{
		authService: authService,
	}
}
