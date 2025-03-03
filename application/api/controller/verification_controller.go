package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	"github.com/purplior/edi-adam/domain/shared/model"
	domain "github.com/purplior/edi-adam/domain/verification"
)

type (
	VerificationController interface {
		Controller
		// 인증코드 요청
		Request() common.Route
		// 인증코드 인증
		Verify() common.Route
	}
)

type (
	verificationController struct {
		verificationService domain.VerificationService
	}
)

func (c *verificationController) GroupPath() string {
	return "/verifications"
}

func (c *verificationController) Routes() []common.Route {
	return []common.Route{
		c.Request(),
		c.Verify(),
	}
}

func (c *verificationController) Request() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/o/fst/request",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			Method string `body:"method,required"`
			Target string `body:"target,required"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		m, err := c.verificationService.Request(
			ctx.Session(),
			model.VerificationMethod(dto.Method),
			dto.Target,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendCreated(struct {
			ID uint `json:"id"`
		}{
			ID: m.ID,
		})
	}

	return route
}

func (c *verificationController) Verify() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/o/fst/verify",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			ID   uint   `body:"id,required"`
			Code string `body:"code,required"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		err := c.verificationService.Verify(
			ctx.Session(),
			dto.ID,
			dto.Code,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(nil)
	}

	return route
}

func NewVerificationController(
	verificationService domain.VerificationService,
) VerificationController {
	return &verificationController{
		verificationService: verificationService,
	}
}
