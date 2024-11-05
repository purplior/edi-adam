package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/auth"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/user"
)

type (
	AuthController interface {
		/**
		 * 이메일로 로그인
		 */
		SignIn_ByEmailVerification() api.HandlerFunc

		/**
		 * 이메일로 로그인 (어드민용)
		 */
		SignIn_ByEmailVerification_ForAdmin() api.HandlerFunc

		/**
		 * 이메일로 회원가입
		 */
		SignUp_ByEmailVerification() api.HandlerFunc

		/**
		 * 토큰 재발급
		 */
		RefreshIdentityToken() api.HandlerFunc
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
		var dto domain.SignInByEmailVerificationRequest
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

func (c *authController) SignIn_ByEmailVerification_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.SignInByEmailVerificationRequest
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
		var dto domain.SignUpByEmailVerificationRequest
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

func NewAuthController(
	authService domain.AuthService,
	cm inner.ContextManager,
) AuthController {
	return &authController{
		authService: authService,
		cm:          cm,
	}
}
