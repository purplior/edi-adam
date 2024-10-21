package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/shared/context"
)

type (
	AssisterFormController interface {
		RegisterOne() api.HandlerFunc
	}
)

type (
	assisterFormController struct {
		assisterFormService domain.AssisterFormService
	}
)

func (c *assisterFormController) RegisterOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.AssisterFormRegisterRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		assisterForm, err := c.assisterFormService.RegisterOne(
			apiCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				AssisterForm domain.AssisterForm `json:"assisterForm"`
			}{
				AssisterForm: assisterForm,
			},
		})
	}
}

func NewAssisterFormController(
	assisterFormService domain.AssisterFormService,
) AssisterFormController {
	return &assisterFormController{
		assisterFormService: assisterFormService,
	}
}
