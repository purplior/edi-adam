package app

import (
	"github.com/purplior/sbec/application/api"
	"github.com/purplior/sbec/application/response"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/user"
)

type (
	UserController interface {
		CheckNicknameExistence() api.HandlerFunc
	}
)

type (
	userController struct {
		userService user.UserService
		cm          inner.ContextManager
	}
)

func (c *userController) CheckNicknameExistence() api.HandlerFunc {
	return func(ctx *api.Context) error {
		var dto struct {
			Nickname string `json:"nickname"`
		}

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		exist, err := c.userService.CheckNicknameExistence(
			innerCtx,
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
}

func NewUserController(
	userService user.UserService,
	cm inner.ContextManager,
) UserController {
	return &userController{
		userService: userService,
		cm:          cm,
	}
}
