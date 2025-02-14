package app

import (
	"github.com/purplior/sbec/application/api"
	"github.com/purplior/sbec/application/response"
	domain "github.com/purplior/sbec/domain/customervoice"
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
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

		request.UserID = ctx.Identity.ID

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		_, err := c.customerVoiceService.RegisterOne(
			innerCtx,
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
