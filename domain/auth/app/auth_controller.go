package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/auth"
	"github.com/purplior/podoroot/domain/shared/inner"
)

type (
	AuthController interface {
		/**
		 * 휴대폰 번호로 로그인
		 */
		SignIn() api.HandlerFunc

		/**
		 * 휴대폰 번호로 회원가입
		 */
		SignUp() api.HandlerFunc

		/**
		 * 토큰 재발급
		 */
		RefreshIdentityToken() api.HandlerFunc

		/**
		 * 휴대폰 번호로 비밀번호 초기화
		 */
		ResetPassword() api.HandlerFunc
	}
)

type (
	authController struct {
		authService domain.AuthService
		cm          inner.ContextManager
	}
)

func (c *authController) SignIn() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.SignInRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		identityToken, identity, err := c.authService.SignIn_ByPhoneNumberVerification(
			innerCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Token    domain.IdentityToken `json:"token"`
				Identity domain.Identity      `json:"identity"`
			}{
				Token:    identityToken,
				Identity: identity,
			},
		})
	}
}

func (c *authController) SignUp() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.SignUpRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		err := c.authService.SignUp_ByPhoneNumberVerification(
			innerCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func (c *authController) RefreshIdentityToken() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.IdentityToken
		if err := ctx.Bind(&dto); err != nil {
			print(err)
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		newIdentityToken, err := c.authService.RefreshIdentityToken(innerCtx, dto)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: newIdentityToken,
		})
	}
}

func (c *authController) ResetPassword() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.ResetPasswordRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		err := c.authService.ResetPassword_ByPhoneNumberVerification(
			innerCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func NewAuthController(
	authService domain.AuthService,
	cm inner.ContextManager,
) AuthController {
	return &authController{
		authService: authService,
		cm:          cm,
	}
}
