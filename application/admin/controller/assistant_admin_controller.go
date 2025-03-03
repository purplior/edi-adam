package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	"github.com/purplior/edi-adam/domain/assistant"
)

type (
	AssistantAdminController interface {
		Controller
		Approve() common.Route
	}
)

type (
	assistantAdminController struct {
		assistantService assistant.AssistantService
	}
)

func (c *assistantAdminController) GroupPath() string {
	return "/assistants"
}

func (c *assistantAdminController) Routes() []common.Route {
	return []common.Route{
		c.Approve(),
	}
}

func (c *assistantAdminController) Approve() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/d/:id/approve",
		Option: common.RouteOption{Admin: true},
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			AssistantID uint     `param:"id"`
			Cost        int      `body:"cost"`
			MetaTags    []string `body:"metaTags"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return err
		}

		if err := c.assistantService.Approve(
			ctx.Session(),
			assistant.ApproveDTO{
				ID: dto.AssistantID,
			},
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(nil)
	}

	return route
}

func NewAssistantAdminController(
	assistantService assistant.AssistantService,
) AssistantAdminController {
	return &assistantAdminController{
		assistantService: assistantService,
	}
}
