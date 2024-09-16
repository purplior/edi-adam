package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	"github.com/podossaem/podoroot/domain/context"
	domain "github.com/podossaem/podoroot/domain/verification"
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

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		emailVerification, err := c.emailVerificationService.RequestCode(
			apiCtx,
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

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		emailVerification, err := c.emailVerificationService.VerifyCode(apiCtx, dto.Email, dto.Code)
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
) EmailVerificationController {
	return &emailVerificationController{
		emailVerificationService: emailVerificationService,
	}
}
