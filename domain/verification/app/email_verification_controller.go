package app

import (
	"github.com/podossaem/root/application/api/controller"
	"github.com/podossaem/root/application/api/response"
	"github.com/podossaem/root/domain/context"
	"github.com/podossaem/root/domain/verification"
)

type (
	EmailVerificationController interface {
		RequestCode() controller.HandlerFunc
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

func NewEmailVerificationController(
	service verification.EmailVerificationService,
) EmailVerificationController {
	return &emailVerificationController{
		service: service,
	}
}
