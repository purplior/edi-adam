package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/logger"
	"github.com/purplior/podoroot/infra/port/podoopenai"
)

type (
	AssisterController interface {
		/**
		 * 어시 정보 가져오기
		 */
		GetInfoOne() api.HandlerFunc

		/**
		 * 어시 실행하기
		 */
		Execute() api.HandlerFunc

		/**
		* 어시 Stream으로 실행하기
		 */
		ExecuteAsStream() api.HandlerFunc
	}
)

type (
	assisterController struct {
		assisterService domain.AssisterService
		execSem         chan struct{}
		cm              inner.ContextManager
	}
)

func (c *assisterController) GetInfoOne() api.HandlerFunc {
	return func(ctx *api.Context) error {
		id := ctx.QueryParam("id")

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		var _assister domain.Assister
		var err error = exception.ErrNotFound

		if len(id) > 0 {
			_assister, err = c.assisterService.GetOne_ByID(
				innerCtx,
				id,
			)
			if err != nil {
				return ctx.SendError(err)
			}
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				AssisterInfo domain.AssisterInfo `json:"assisterInfo"`
			}{
				AssisterInfo: _assister.ToInfo(),
			},
		})
	}
}

func (c *assisterController) Execute() api.HandlerFunc {
	return func(ctx *api.Context) error {
		userID := ""
		if ctx.Identity != nil {
			userID = ctx.Identity.ID
		}

		id := ctx.QueryParam("id")
		if len(id) == 0 {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		var dto struct {
			Inputs []domain.AssisterInput `json:"inputs"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		innerCtx, cancel := c.cm.NewStreamingContext()
		defer cancel()

		result, err := c.assisterService.Request(
			innerCtx,
			userID,
			id,
			dto.Inputs,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Result string `json:"result"`
			}{
				Result: result,
			},
		})
	}
}

func (c *assisterController) ExecuteAsStream() api.HandlerFunc {
	return func(ctx *api.Context) error {
		userID := ""
		if ctx.Identity != nil {
			userID = ctx.Identity.ID
		}

		assisterID := ctx.QueryParam("aid")
		if len(assisterID) == 0 {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		var dto struct {
			Inputs []domain.AssisterInput `json:"inputs"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		select {
		case c.execSem <- struct{}{}: // 세마포어를 사용해서 리소스 소비
			defer func() { <-c.execSem }()

			innerCtx, cancel := c.cm.NewStreamingContext()
			defer cancel()

			err := c.assisterService.RequestAsStream(
				innerCtx,
				userID,
				assisterID,
				dto.Inputs,
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
				case podoopenai.ErrOnStream:
					return nil
				case exception.ErrBadRequest:
					return ctx.String(http.StatusBadRequest, "입력폼을 다시 확인해주세요")
				case exception.ErrNoPodo:
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
}

func NewAssisterController(
	assisterService domain.AssisterService,
	cm inner.ContextManager,
) AssisterController {
	return &assisterController{
		assisterService: assisterService,
		execSem:         make(chan struct{}, 500), // 동시 요청 500개로 제한
		cm:              cm,
	}
}
