package app

import (
	"github.com/podossaem/podoroot/application/api/controller"
	"github.com/podossaem/podoroot/application/api/response"
	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/domain/user"
)

type (
	UserController interface {
		/**
		 * 이메일로 회원가입
		 */
		SignUpByEmailVerification() controller.HandlerFunc
	}

	userController struct {
		userService user.UserService
	}
)

func (c *userController) SignUpByEmailVerification() controller.HandlerFunc {
	return func(ctx *controller.Context) error {
		var dto user.SignUpRequest

		if err := ctx.Bind(&dto); err != nil {
			return ctx.SendError(err)
		}

		apiCtx, cancel := context.NewAPIContext()
		defer cancel()

		signedUpUser, err := c.userService.SignUpByEmailVerification(
			apiCtx,
			dto,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				User user.User `json:"user"`
			}{
				User: signedUpUser,
			},
		})
	}
}

func NewUserController(
	userService user.UserService,
) UserController {
	return &userController{
		userService: userService,
	}
}
