package app

import (
	"github.com/podossaem/podoroot/application/api/controller"
	"github.com/podossaem/podoroot/application/api/response"
	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/domain/exception"
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
		service verification.EmailVerificationService
	}
)

func (c *emailVerificationController) RequestCode() controller.HandlerFunc {
	return func(ctx *controller.Context) error {
		var dto struct {
			Email string `json:"email"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(response.ErrorResponse{
				Status:  response.Status_InternalServerError,
				Message: response.Message_ErrorNormal,
			})
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		if _, err := c.service.RequestCode(apiCtx, dto.Email); err != nil {
			return ctx.SendError(response.ErrorResponse{
				Status:  response.Status_InternalServerError,
				Message: response.Message_ErrorNormal,
			})
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
			return ctx.SendError(response.ErrorResponse{
				Status:  response.Status_InternalServerError,
				Message: response.Message_ErrorNormal,
			})
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		emailVerification, err := c.service.VerifyCode(apiCtx, dto.Email, dto.Code)
		if err != nil {
			status := response.Status_InternalServerError

			switch err {
			case verification.ErrInvalidCode:
				status = response.Status_BadRequest
			case verification.ErrAlreadyVerified:
				status = response.Status_BadRequest
			case exception.ErrNoDocuments:
				status = response.Status_BadRequest
			}

			return ctx.SendError(response.ErrorResponse{
				Status:  status,
				Message: err.Error(),
			})
		}

		responseData := struct {
			ID string `json:"id"`
		}{
			ID: emailVerification.ID,
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Ok,
			Data:   responseData,
		})
	}
}

func NewEmailVerificationController(
	service verification.EmailVerificationService,
) EmailVerificationController {
	return &emailVerificationController{
		service: service,
	}
}
