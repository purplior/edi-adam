package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api/response"
	"github.com/podossaem/podoroot/domain/exception"
)

type (
	Context struct {
		echo.Context
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

	switch err {
	case exception.ErrBadRequest:
		status = response.Status_BadRequest
	case exception.ErrAlreadyConsumed:
		status = response.Status_BadRequest
	case exception.ErrInvalidVerificationCode:
		status = response.Status_BadRequest
	case exception.ErrAlreadyVerified:
		status = response.Status_BadRequest
	}

	return ctx.JSON(status, response.ErrorResponse{
		Status:  status,
		Message: err.Error(),
	})
}
