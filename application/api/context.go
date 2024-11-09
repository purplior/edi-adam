package api

import (
	"log"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/application/response"
	"github.com/podossaem/podoroot/domain/auth"
	"github.com/podossaem/podoroot/domain/shared/constant"
	"github.com/podossaem/podoroot/domain/shared/exception"
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
	message := "internal server error"

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
	}

	switch status {
	case response.Status_BadRequest:
		message = err.Error()
	case response.Status_NotAcceptable:
		message = err.Error()
	case response.Status_Unauthorized:
		message = err.Error()
	case response.Status_InternalServerError:
		log.Println(err.Error())
	}

	if config.Phase() == constant.Phase_Local {
		log.Println(err)
		log.Printf("Error: %v\nStack Trace:\n%s", err, debug.Stack())
	}

	return ctx.JSON(status, response.ErrorResponse{
		Status:  status,
		Message: message,
	})
}
