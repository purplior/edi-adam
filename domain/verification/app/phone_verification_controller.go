package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	"github.com/purplior/podoroot/domain/shared/inner"
	domain "github.com/purplior/podoroot/domain/verification"
)

type (
	PhoneVerificationController interface {
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
	phoneVerificationController struct {
		phoneVerificationService domain.PhoneVerificationService
		cm                       inner.ContextManager
	}
)

func (c *phoneVerificationController) RequestCode() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			PhoneNumber string `json:"phoneNumber"`
		}
		isTestMode := len(ctx.QueryParam("test")) > 0

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		verification, err := c.phoneVerificationService.RequestCode(
			innerCtx,
			dto.PhoneNumber,
			isTestMode,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		responseData := struct {
			ID string `json:"id"`
		}{
			ID: verification.ID,
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Created,
			Data:   responseData,
		})
	}
}

func (c *phoneVerificationController) VerifyCode() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			PhoneNumber string `json:"phoneNumber"`
			Code        string `json:"code"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		verification, err := c.phoneVerificationService.VerifyCode(innerCtx, dto.PhoneNumber, dto.Code)
		if err != nil {
			return ctx.SendError(err)
		}

		responseData := struct {
			ID string `json:"id"`
		}{
			ID: verification.ID,
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: responseData,
		})
	}
}

func NewPhoneVerificationController(
	phoneVerificationService domain.PhoneVerificationService,
	cm inner.ContextManager,
) PhoneVerificationController {
	return &phoneVerificationController{
		phoneVerificationService: phoneVerificationService,
		cm:                       cm,
	}
}
