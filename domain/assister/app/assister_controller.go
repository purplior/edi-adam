package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/response"
	domain "github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/shared/pagination"
	"github.com/podossaem/podoroot/infra/port/podoopenai"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	AssisterController interface {
		/**
		 * 쌤비서 실행하기
		 */
		ExecuteAsResult() api.HandlerFunc

		/**
		 * 쌤비서 Stream으로 실행하기
		 */
		ExecuteAsStream() api.HandlerFunc

		/**
		 * 쌤비서 실행기 가져오기
		 */
		GetOne_ForAdmin() api.HandlerFunc

		/**
		 * 쌤비서 실행기 목록 가져오기
		 */
		GetPaginatedList_ForAdmin() api.HandlerFunc

		/**
		 * 쌤비서 실행기 수정하기
		 */
		PutOne_ForAdmin() api.HandlerFunc

		/**
		 * 쌤비서 실행기 생성하기 (어드민용)
		 */
		CreateOne_ForAdmin() api.HandlerFunc
	}
)

type (
	assisterController struct {
		execSem         chan struct{}
		assisterService domain.AssisterService
		cm              inner.ContextManager
	}
)

func (c *assisterController) ExecuteAsResult() api.HandlerFunc {
	return func(ctx *api.Context) error {
		userId := ""
		if ctx.Identity != nil {
			userId = ctx.Identity.ID
		}

		assisterID := ctx.QueryParam("aid")
		if len(assisterID) == 0 {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		var dto struct {
			Inputs []assisterform.AssisterInput `json:"inputs"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		innerCtx, cancel := c.cm.NewStreamingContext()
		defer cancel()

		result, err := c.assisterService.Request(
			innerCtx,
			userId,
			assisterID,
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
		userId := ""
		if ctx.Identity != nil {
			userId = ctx.Identity.ID
		}

		assisterID := ctx.QueryParam("aid")
		if len(assisterID) == 0 {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		var dto struct {
			Inputs []assisterform.AssisterInput `json:"inputs"`
		}
		if err := ctx.Bind(&dto); err != nil {
			return ctx.String(http.StatusBadRequest, exception.ErrBadRequest.Error())
		}

		select {
		case c.execSem <- struct{}{}: // 세마포어를 사용해서 리소스 소비
			defer func() { <-c.execSem }()

			innerCtx, cancel := c.cm.NewStreamingContext()
			defer cancel()

			err := c.assisterService.RequestStream(
				innerCtx,
				userId,
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

func (c *assisterController) GetOne_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		assisterID := ctx.QueryParam("id")

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assister, err := c.assisterService.GetOne_ByID(innerCtx, assisterID)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Assister domain.Assister `json:"assister"`
			}{
				Assister: assister,
			},
		})
	}
}

func (c *assisterController) GetPaginatedList_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		assistantID := ctx.QueryParam("assistant_id")
		page := dt.Int(ctx.QueryParam("p"))
		pageSize := dt.Int(ctx.QueryParam("ps"))

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		assisters, meta, err := c.assisterService.GetPaginatedList_ByAssistant(
			innerCtx,
			assistantID,
			page,
			pageSize,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Assisters []domain.Assister         `json:"assisters"`
				Meta      pagination.PaginationMeta `json:"meta"`
			}{
				Assisters: assisters,
				Meta:      meta,
			},
		})
	}
}

func (c *assisterController) PutOne_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			Assister domain.Assister `json:"assister"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		if err := c.assisterService.PutOne(
			innerCtx,
			dto.Assister,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{})
	}
}

func (c *assisterController) CreateOne_ForAdmin() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			Assister domain.Assister `json:"assister"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		if err := c.assisterService.CreateOne(
			innerCtx,
			dto.Assister,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Status: response.Status_Created,
		})
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
