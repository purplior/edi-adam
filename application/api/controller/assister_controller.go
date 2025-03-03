package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/purplior/edi-adam/application/common"
	domain "github.com/purplior/edi-adam/domain/assister"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/logger"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/port/openai"
)

type (
	AssisterController interface {
		Controller
		// 실행을 위한 어시 정보 가져오기
		Get_Excutable() common.Route
		// 어시 실행하기
		Execute() common.Route
		// 어시 Stream으로 실행하기
		ExecuteAsStream() common.Route
	}
)

type (
	assisterController struct {
		assisterService domain.AssisterService
		execSem         chan struct{}
	}
)

func (c *assisterController) GroupPath() string {
	return "/assisters"
}

func (c *assisterController) Routes() []common.Route {
	return []common.Route{
		c.Get_Excutable(),
		c.Execute(),
		c.ExecuteAsStream(),
	}
}

func (c *assisterController) Get_Excutable() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/o/:id/excutable",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			ID string `param:"id,required"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		m, err := c.assisterService.Get(
			ctx.Session(),
			domain.QueryOption{
				ID: dto.ID,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			ExcutableAssister model.ExcutableAssister `json:"excutableAssister"`
		}{
			ExcutableAssister: m.ToExcutable(),
		})
	}

	return route
}

func (c *assisterController) Execute() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/m/:id/plaintext",
		Option: common.RouteOption{Member: true, Timeout: time.Duration(5) * time.Minute},
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			ID     string                `param:"id,required"`
			Inputs []model.AssisterInput `body:"inputs"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		result, err := c.assisterService.Request(
			ctx.Session(),
			domain.RequestDTO{
				ID:     dto.ID,
				UserID: ctx.Identity.ID,
				Inputs: dto.Inputs,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Result string `json:"result"`
		}{
			Result: result,
		})
	}

	return route
}

func (c *assisterController) ExecuteAsStream() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/m/:id/stream",
		Option: common.RouteOption{Member: true, Timeout: time.Duration(5) * time.Minute},
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			ID     string                `param:"id,required"`
			Inputs []model.AssisterInput `body:"inputs"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		select {
		case c.execSem <- struct{}{}: // 세마포어를 사용해서 리소스 소비
			defer func() { <-c.execSem }()

			err := c.assisterService.RequestAsStream(
				ctx.Session(),
				domain.RequestDTO{
					ID:     dto.ID,
					UserID: ctx.Identity.ID,
					Inputs: dto.Inputs,
				},
				func() error {
					ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
					ctx.Response().WriteHeader(http.StatusOK)

					if f, ok := ctx.Response().Writer.(http.Flusher); ok {
						f.Flush()
					}

					return nil
				},
				func(msg string) error {
					_, err := ctx.Response().Writer.Write([]byte(msg))
					if err != nil {
						return err
					}

					if f, ok := ctx.Response().Writer.(http.Flusher); ok {
						f.Flush()
					}

					return nil
				},
			)

			if err != nil {
				logger.Debug("%s", err.Error())

				switch err {
				case openai.ErrOnStream:
					return nil
				case exception.ErrBadRequest:
					return ctx.String(http.StatusBadRequest, "입력폼을 다시 확인해주세요")
				case exception.ErrNoCoin:
					return ctx.String(http.StatusForbidden, err.Error())
				default:
					return ctx.String(http.StatusInternalServerError, "일시적인 서버 오류가 발생했어요")
				}
			}

			return nil
		default:
			return ctx.String(http.StatusTooManyRequests, "현재 이용자가 매우 많아요, 잠시 후 다시 시도해주세요.")
		}
	}

	return route
}

func NewAssisterController(
	assisterService domain.AssisterService,
) AssisterController {
	return &assisterController{
		assisterService: assisterService,
		execSem:         make(chan struct{}, 500), // 동시 요청 500개로 제한
	}
}
