package app

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
	domain "github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/infra/port/podoopenai"
)

type (
	AssisterController interface {
		/**
		 * 쌤비서 실행하기
		 */
		Execute() api.HandlerFunc
	}
)

type (
	assisterController struct {
		execSem         chan struct{}
		assisterService domain.AssisterService
	}
)

func (c *assisterController) Execute() api.HandlerFunc {
	return func(ctx *api.Context) error {
		assisterID := ctx.Param("assister_id")
		if len(assisterID) == 0 {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		var dto struct {
			Inputs map[string]interface{} `json:"inputs"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		select {
		case c.execSem <- struct{}{}: // 세마포어를 사용해서 리소스 소비
			defer func() { <-c.execSem }()

			apiCtx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Minute)*5)
			defer cancel()

			err := c.assisterService.RequestStream(
				apiCtx,
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

			if err != nil && err != podoopenai.ErrOnStream {
				return nil
			}

			return nil
		default:
			return ctx.String(http.StatusTooManyRequests, "현재 이용자가 매우 많아요, 잠시 후 다시 시도해주세요.")
		}
	}
}

func NewAssisterController(
	assisterService domain.AssisterService,
) AssisterController {
	return &assisterController{
		assisterService: assisterService,
		execSem:         make(chan struct{}, 500), // 동시 요청 500개로 제한
	}
}
