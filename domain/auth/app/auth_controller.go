package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/auth"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/user"
)

type (
	AuthController interface {
		/**
		 * 이메일로 로그인
		 */
		SignIn_ByEmailVerification() api.HandlerFunc

		/**
		 * 휴대폰 번호로 로그인
		 */
		SignIn_ByPhoneNumberVerification() api.HandlerFunc

		/**
		 * 이메일로 로그인 (어드민용)
		 */
		SignIn_ByEmailVerification_ForAdmin() api.HandlerFunc

		/**
		 * 이메일로 회원가입
		 */
		SignUp_ByEmailVerification() api.HandlerFunc

		/**
		 * 휴대폰 번호로 회원가입
		 */
		SignUp_ByPhoneNumberVerification() api.HandlerFunc

		/**
		 * 토큰 재발급
		 */
		RefreshIdentityToken() api.HandlerFunc

		/**
		 * 휴대폰 번호로 비밀번호 초기화
		 */
		ResetPassword_ByPhoneNumberVerification() api.HandlerFunc
	}
)

type (
	authController struct {
		authService domain.AuthService
		cm          inner.ContextManager
	}
)

func (c *authController) SignIn_ByEmailVerification() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.SignInRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		identityToken, identity, err := c.authService.SignIn_ByEmailVerification(
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

func (c *authController) SignIn_ByPhoneNumberVerification() api.HandlerFunc {
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

func (c *authController) SignIn_ByEmailVerification_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.SignInRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		identityToken, identity, err := c.authService.SignIn_ByEmailVerification(
			innerCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}
		if identity.Role != user.Role_Master {
			return ctx.SendError(exception.ErrUnauthorized)
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

func (c *authController) SignUp_ByEmailVerification() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.SignUpRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		err := c.authService.SignUp_ByEmailVerification(
			innerCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func (c *authController) SignUp_ByPhoneNumberVerification() api.HandlerFunc {
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

func (c *authController) ResetPassword_ByPhoneNumberVerification() api.HandlerFunc {
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
