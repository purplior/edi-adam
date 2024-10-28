package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
)

type (
	AssisterFormController interface {
		/**
		 * 쌤비서 폼 등록하기
		 */
		RegisterOne() api.HandlerFunc

		/**
		 * 쌤비서 폼 가져오기
		 */
		GetOne() api.HandlerFunc
	}
)

type (
	assisterFormController struct {
		assisterFormService domain.AssisterFormService
		cm                  inner.ContextManager
	}
)

func (c *assisterFormController) RegisterOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto domain.AssisterFormRegisterRequest
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assisterForm, err := c.assisterFormService.RegisterOne(
			innerCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Created,
			Data: struct {
				AssisterForm domain.AssisterForm `json:"assisterForm"`
			}{
				AssisterForm: assisterForm,
			},
		})
	}
}

func (c *assisterFormController) GetOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		assisterID := ctx.QueryParam("assister_id")

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		var assisterForm domain.AssisterForm
		var err error = exception.ErrNotFound

		if len(assisterID) > 0 {
			assisterForm, err = c.assisterFormService.GetOneByAssisterID(
				innerCtx,
				assisterID,
			)
			if err != nil {
				return ctx.SendError(err)
			}
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
	cm inner.ContextManager,
) AssisterFormController {
	return &assisterFormController{
		assisterFormService: assisterFormService,
		cm:                  cm,
	}
}
