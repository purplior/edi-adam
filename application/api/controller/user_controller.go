package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	"github.com/purplior/edi-adam/application/response"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/domain/user"
	"github.com/purplior/edi-adam/domain/wallet"
)

type (
	UserController interface {
		Controller
		// 내 정보 가져오기
		GetMine_Detail() common.Route
		// 나의 코인 확인하기
		GetMine_Coin() common.Route
		// 닉네임 중복 확인하기
		CheckNicknameExistence() common.Route
	}
)

type (
	userController struct {
		userService   user.UserService
		walletService wallet.WalletService
	}
)

func (c *userController) GroupPath() string {
	return "/users"
}

func (c *userController) Routes() []common.Route {
	return []common.Route{
		c.GetMine_Detail(),
		c.GetMine_Coin(),
		c.CheckNicknameExistence(),
	}
}

func (c *userController) GetMine_Detail() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/m/fst/detail",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var m model.User
		var detail model.UserDetail

		if m, err = c.userService.Get(
			ctx.Session(),
			user.QueryOption{
				ID: ctx.Identity.ID,
			},
		); err != nil {
			return ctx.SendError(err)
		}
		if detail, err = m.ToDetail(); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			UserDetail model.UserDetail `json:"userDetail"`
		}{
			UserDetail: detail,
		})
	}

	return route
}

func (c *userController) GetMine_Coin() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/m/fst/coin",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) error {
		wallet, err := c.walletService.Get(
			ctx.Session(),
			wallet.QueryOption{
				OwnerID: ctx.Identity.ID,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Coin int64 `json:"coin"`
		}{
			Coin: wallet.Coin,
		})
	}

	return route
}

func (c *userController) CheckNicknameExistence() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/o/fst/nickname-check",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			Nickname string `body:"nickname"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		exist, err := c.userService.CheckNicknameExistence(
			ctx.Session(),
			dto.Nickname,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				Exist bool `json:"exist"`
			}{
				Exist: exist,
			},
		})
	}

	return route
}

func NewUserController(
	userService user.UserService,
	walletService wallet.WalletService,
) UserController {
	return &userController{
		userService:   userService,
		walletService: walletService,
	}
}
