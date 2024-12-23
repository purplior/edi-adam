package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/customervoice"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
)

type (
	CustomerVoiceController interface {
		RegisterOne() api.HandlerFunc
	}
)

type (
	customerVoiceController struct {
		customerVoiceService domain.CustomerVoiceService
		cm                   inner.ContextManager
	}
)

func (c *customerVoiceController) RegisterOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		if ctx.Identity == nil {
			return ctx.SendError(exception.ErrBadRequest)
		}

		var request domain.CustomerVoiceRegisterRequest
		if err := ctx.Bind(&request); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		_, err := c.customerVoiceService.RegisterOne(
			innerCtx,
			ctx.Identity.ID,
			request,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Created,
		})
	}
}

func NewCustomerVoiceController(
	customerVoiceService domain.CustomerVoiceService,
	cm inner.ContextManager,
) CustomerVoiceController {
	return &customerVoiceController{
		customerVoiceService: customerVoiceService,
		cm:                   cm,
	}
}
