package api

import (
	"log"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/application/response"
	"github.com/purplior/podoroot/domain/auth"
	"github.com/purplior/podoroot/domain/shared/constant"
	"github.com/purplior/podoroot/domain/shared/exception"
)

type (
	Context struct {
		echo.Context
		Identity *auth.Identity
	}
)

func (ctx Context) SendJSON(jsonResponse response.JSONResponse) error {
	if jsonResponse.Status == 0 {
		jsonResponse.Status = response.Status_Ok
	}

	return ctx.JSON(jsonResponse.Status, jsonResponse)
}

func (ctx Context) SendError(err error) error {
	status := response.Status_InternalServerError
	message := "일시적인 서버 오류가 발생했어요"

	switch err {
	case exception.ErrBadRequest:
		status = response.Status_BadRequest
	case exception.ErrAlreadyConsumed:
		status = response.Status_BadRequest
	case exception.ErrNotConsumed:
		status = response.Status_BadRequest
	case exception.ErrInvalidVerificationCode:
		status = response.Status_BadRequest
	case exception.ErrAlreadyVerified:
		status = response.Status_BadRequest
	case exception.ErrAlreadyReceived:
		status = response.Status_BadRequest
	case exception.ErrNotAcceptable:
		status = response.Status_NotAcceptable
	case exception.ErrAlreadySignedUp:
		status = response.Status_NotAcceptable
	case exception.ErrUnauthorized:
		status = response.Status_Unauthorized
	case exception.ErrNoRecord:
		status = response.Status_NotFound
	case exception.ErrNoSignedUpPhone:
		status = response.Status_Unprocessable
	case exception.ErrPhoneVerificationExceed:
		status = response.Status_Unprocessable
	case exception.ErrNotAllowedNickname:
		status = response.Status_BadRequest
	}

	if status != response.Status_InternalServerError {
		message = err.Error()
	}

	if config.Phase() != constant.Phase_Production {
		log.Println(err.Error())
		log.Printf("Error: %v\nStack Trace:\n%s", err, debug.Stack())
	}

	return ctx.JSON(status, response.ErrorResponse{
		Status:  status,
		Message: message,
	})
}
