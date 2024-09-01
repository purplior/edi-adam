package app

import (
	"github.com/podossaem/podoroot/application/api/controller"
	"github.com/podossaem/podoroot/application/api/response"
	domain "github.com/podossaem/podoroot/domain/auth"
	"github.com/podossaem/podoroot/domain/context"
)

type (
	AuthController interface {
		/**
		 * 이메일로 로그인
		 */
		SignInByEmailVerification() controller.HandlerFunc

		/**
		* 이메일로 회원가입
		 */
		SignUpByEmailVerification() controller.HandlerFunc
	}
)

type (
	authController struct {
		authService domain.AuthService
	}
)

func (c *authController) SignInByEmailVerification() controller.HandlerFunc {
	return func(ctx *controller.Context) error {
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

func (c *authController) SignUpByEmailVerification() controller.HandlerFunc {
	return func(ctx *controller.Context) error {
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

func NewAuthController(
	authService domain.AuthService,
) AuthController {
	return &authController{
		authService: authService,
	}
}
