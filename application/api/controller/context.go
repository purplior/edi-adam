package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/root/application/api/response"
)

type (
	Context struct {
		echo.Context
	}
)

func (ctx Context) SendJSON(response response.JSONResponse) error {
	return ctx.JSON(response.Status, response)
}

func (ctx Context) SendError(response response.ErrorResponse) error {
	return ctx.JSON(response.Status, response)
}
