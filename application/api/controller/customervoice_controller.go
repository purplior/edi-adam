package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	domain "github.com/purplior/edi-adam/domain/customervoice"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	CustomerVoiceController interface {
		Controller
		// 고객문의 등록
		Register() common.Route
	}
)

type (
	customerVoiceController struct {
		customerVoiceService domain.CustomerVoiceService
	}
)

func (c *customerVoiceController) GroupPath() string {
	return "/customervoices"
}

func (c *customerVoiceController) Routes() []common.Route {
	return []common.Route{
		c.Register(),
	}
}

func (c *customerVoiceController) Register() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/m/fst",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			Type    string `body:"type,required"`
			Content string `body:"content,required"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		if _, err = c.customerVoiceService.Register(
			ctx.Session(),
			domain.CustomerVoiceRegisterDTO{
				UserID:  ctx.Identity.ID,
				Type:    model.CustomerVoiceType(dto.Type),
				Content: dto.Content,
			},
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendCreated(nil)
	}

	return route
}

func NewCustomerVoiceController(
	customerVoiceService domain.CustomerVoiceService,
) CustomerVoiceController {
	return &customerVoiceController{
		customerVoiceService: customerVoiceService,
	}
}
