package app

import (
	"github.com/podossaem/podoroot/application/api/controller"
	"github.com/podossaem/podoroot/application/api/response"
	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/domain/verification"
)

type (
	EmailVerificationController interface {
		/**
		 * 인증코드 요청
		 */
		RequestCode() controller.HandlerFunc

		/**
		 * 인증코드 인증
		 */
		VerifyCode() controller.HandlerFunc
	}

	emailVerificationController struct {
		emailVerificationService verification.EmailVerificationService
	}
)

func (c *emailVerificationController) RequestCode() controller.HandlerFunc {
	return func(ctx *controller.Context) error {
		var dto struct {
			Email string `json:"email"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		if _, err := c.emailVerificationService.RequestCode(apiCtx, dto.Email); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Created,
			Data:   nil,
		})
	}
}

func (c *emailVerificationController) VerifyCode() controller.HandlerFunc {
	return func(ctx *controller.Context) error {
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
	emailVerificationService verification.EmailVerificationService,
) EmailVerificationController {
	return &emailVerificationController{
		emailVerificationService: emailVerificationService,
	}
}
