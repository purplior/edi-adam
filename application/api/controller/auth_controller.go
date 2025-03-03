package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	domain "github.com/purplior/edi-adam/domain/auth"
	"github.com/purplior/edi-adam/domain/shared/inner"
)

type (
	AuthController interface {
		Controller
		// 내 정보(인증에 포함된 간단한 식별정보) 가져오기
		GetMine_Identity() common.Route
		// 나의 임시 액세스 토큰 발급하기 (유효기간 1시간)
		GetMine_TempAccessToken() common.Route
		// 휴대폰 번호로 로그인
		SignIn() common.Route
		// 휴대폰 번호로 회원가입
		SignUp() common.Route
	}
)

type (
	authController struct {
		authService domain.AuthService
	}
)

func (c *authController) GroupPath() string {
	return "/auth"
}

func (c *authController) Routes() []common.Route {
	return []common.Route{
		c.GetMine_Identity(),
		c.GetMine_TempAccessToken(),
		c.SignIn(),
		c.SignUp(),
	}
}

func (c *authController) GetMine_Identity() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/m/get_identity",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) error {
		return ctx.SendOK(ctx.Identity)
	}

	return route
}

func (c *authController) GetMine_TempAccessToken() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/m/get_tmpat",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) error {
		accessToken, err := c.authService.GetTempAccessToken(
			ctx.Session(),
			*ctx.Identity,
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			TempAccessToken string `json:"tempAccessToken"`
		}{
			TempAccessToken: accessToken,
		})
	}

	return route
}

func (c *authController) SignIn() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/o/signin",
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			VerificationID uint `body:"vid,required"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		identityToken, identity, err := c.authService.SignIn(
			ctx.Session(),
			domain.SignInDTO{
				VerificationID: dto.VerificationID,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Token    domain.IdentityToken `json:"token"`
			Identity inner.Identity       `json:"identity"`
		}{
			Token:    identityToken,
			Identity: identity,
		})
	}

	return route
}

func (c *authController) SignUp() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/o/signup",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			VerificationID   uint   `body:"vid,required"`
			Nickname         string `body:"nickname,required"`
			IsMarketingAgree bool   `body:"isMarketingAgree,required"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		identityToken, identity, err := c.authService.SignUp(
			ctx.Session(),
			domain.SignUpDTO{
				VerificationID:   dto.VerificationID,
				Nickname:         dto.Nickname,
				IsMarketingAgree: dto.IsMarketingAgree,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Token    domain.IdentityToken `json:"token"`
			Identity inner.Identity       `json:"identity"`
		}{
			Token:    identityToken,
			Identity: identity,
		})
	}

	return route
}

func NewAuthController(
	authService domain.AuthService,
) AuthController {
	return &authController{
		authService: authService,
	}
}
