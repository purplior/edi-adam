package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/me"
)

type (
	MeController interface {
		/**
		 * 내 정보 가져오기
		 */
		GetMyIdentity() api.HandlerFunc
	}
)

type (
	meController struct {
		meService domain.MeService
	}
)

func (c *meController) GetMyIdentity() api.HandlerFunc {
	return func(ctx *api.Context) error {
		identity := ctx.Identity

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Ok,
			Data:   identity,
		})
	}
}

func NewMeController(
	meService domain.MeService,
) MeController {
	return &meController{
		meService: meService,
	}
}
