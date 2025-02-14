package app

import (
	"github.com/purplior/sbec/application/api"
	"github.com/purplior/sbec/application/response"
	"github.com/purplior/sbec/domain/shared/inner"
	domain "github.com/purplior/sbec/domain/verification"
)

type (
	EmailVerificationController interface {
		/**
		 * 인증코드 요청
		 */
		RequestCode() api.HandlerFunc

		/**
		 * 인증코드 인증
		 */
		VerifyCode() api.HandlerFunc
	}
)

type (
	emailVerificationController struct {
		emailVerificationService domain.EmailVerificationService
		cm                       inner.ContextManager
	}
)

func (c *emailVerificationController) RequestCode() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			Email string `json:"email"`
		}
		isTestMode := len(ctx.QueryParam("test")) > 0

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		emailVerification, err := c.emailVerificationService.RequestCode(
			innerCtx,
			dto.Email,
			isTestMode,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		responseData := struct {
			ID string `json:"id"`
		}{
			ID: emailVerification.ID,
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Created,
			Data:   responseData,
		})
	}
}

func (c *emailVerificationController) VerifyCode() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		emailVerification, err := c.emailVerificationService.VerifyCode(innerCtx, dto.Email, dto.Code)
		if err != nil {
			return ctx.SendError(err)
		}

		responseData := struct {
			ID string `json:"id"`
		}{
			ID: emailVerification.ID,
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: responseData,
		})
	}
}

func NewEmailVerificationController(
	emailVerificationService domain.EmailVerificationService,
	cm inner.ContextManager,
) EmailVerificationController {
	return &emailVerificationController{
		emailVerificationService: emailVerificationService,
		cm:                       cm,
	}
}
