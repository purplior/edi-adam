package app

import (
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/challenge"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	ChallengeController interface {
		GetPaginatedList() api.HandlerFunc
	}
)

type (
	challengeController struct {
		challengeService domain.ChallengeService
		cm               inner.ContextManager
	}
)

func (c *challengeController) GetPaginatedList() api.HandlerFunc {
	return func(ctx *api.Context) error {
		userID := ctx.Identity.ID

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		challenges, err := c.challengeService.GetPaginatedListByUserID(
			innerCtx,
			userID,
			dt.Int(ctx.QueryParam("psize")),
			dt.Int(ctx.QueryParam("p")),
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Challenges []domain.Challenge `json:"challenges"`
			}{
				Challenges: challenges,
			},
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
