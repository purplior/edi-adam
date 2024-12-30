package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/user"
	domain "github.com/purplior/podoroot/domain/verification"
)

type (
	PhoneVerificationController interface {
		/**
		 * 인증코드 요청
		 */
		RequestCode() api.HandlerFunc

		/**
		 * 가입한 유저의 인증코드 요청
		 */
		RequestCodeOfJoinedUser() api.HandlerFunc

		/**
		 * 인증코드 인증
		 */
		VerifyCode() api.HandlerFunc
	}
)

type (
	phoneVerificationController struct {
		phoneVerificationService domain.PhoneVerificationService
		userService              user.UserService
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

func (c *phoneVerificationController) RequestCodeOfJoinedUser() api.HandlerFunc {
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

		_user, err := c.userService.GetOne_ByAccount(
			innerCtx,
			user.JoinMethod_PhoneNumber,
			dto.PhoneNumber,
		)
		if len(_user.ID) == 0 || err != nil {
			return ctx.SendError(exception.ErrNoSignedUpPhone)
		}

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
	userService user.UserService,
	cm inner.ContextManager,
) PhoneVerificationController {
	return &phoneVerificationController{
		phoneVerificationService: phoneVerificationService,
		userService:              userService,
		cm:                       cm,
	}
}
