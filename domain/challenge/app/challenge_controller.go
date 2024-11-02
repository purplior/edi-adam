package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/challenge"
	"github.com/podossaem/podoroot/domain/shared/inner"
)

type (
	ChallengeController interface {
		ReceiveOne() api.HandlerFunc
	}
)

type (
	challengeController struct {
		challengeService domain.ChallengeService
		cm               inner.ContextManager
	}
)

func (c challengeController) ReceiveOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		userID := ctx.Identity.ID

		var dto struct {
			ID string `json:"id"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		if err := c.challengeService.ReceiveOne(
			innerCtx,
			dto.ID,
			userID,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: nil,
		})
	}
}

func NewChallengeController(
	challengeService domain.ChallengeService,
	cm inner.ContextManager,
) ChallengeController {
	return &challengeController{
		challengeService: challengeService,
		cm:               cm,
	}
}
