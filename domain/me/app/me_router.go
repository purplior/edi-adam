package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
)

type (
	MeRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	meRouter struct {
		meController MeController
	}
)

func (r *meRouter) Attach(router *echo.Group) {
	rg := router.Group("/me")

	/**
	 * 사용자 정보 관련
	 */
	rg.GET(
		"/identity",
		api.Handler(
			r.meController.GetMyIdentity(),
			api.HandlerFuncOption{},
		),
	)
	rg.GET(
		"/detail",
		api.Handler(
			r.meController.GetMyDetail(),
			api.HandlerFuncOption{},
		),
	)
	rg.GET(
		"/podo",
		api.Handler(
			r.meController.GetMyPodo(),
			api.HandlerFuncOption{},
		),
	)
	rg.GET(
		"/temp/at",
		api.Handler(
			r.meController.GetTempAccessToken(),
			api.HandlerFuncOption{},
		),
	)

	/**
	 * 어시 관련
	 */
	rg.GET(
		"/assistant/one",
		api.Handler(
			r.meController.GetMyAssistantOne(),
			api.HandlerFuncOption{},
		),
	)
	rg.GET(
		"/assistant/detail-one",
		api.Handler(
			r.meController.GetMyAssistantDetailOne(),
			api.HandlerFuncOption{},
		),
	)
	rg.GET(
		"/assistant/info-plist",
		api.Handler(
			r.meController.GetMyAssistantInfoPaginatedList(),
			api.HandlerFuncOption{},
		),
	)
	rg.POST(
		"/assistant",
		api.Handler(
			r.meController.RegisterMyAssistantOne(),
			api.HandlerFuncOption{},
		),
	)
	rg.PATCH(
		"/assistant/:id",
		api.Handler(
			r.meController.UpdateMyAssistantOne(),
			api.HandlerFuncOption{},
		),
	)
	rg.DELETE(
		"/assistant/:id",
		api.Handler(
			r.meController.RemoveMyAssistantOne(),
			api.HandlerFuncOption{},
		),
	)

	/**
	 * 북마크 관련
	 */
	rg.GET(
		"/bookmark/one",
		api.Handler(
			r.meController.GetMyBookmarkOne(),
			api.HandlerFuncOption{},
		),
	)
	rg.GET(
		"/bookmark/plist",
		api.Handler(
			r.meController.GetMyBookmarkPaginatedList(),
			api.HandlerFuncOption{},
		),
	)
	rg.POST(
		"/bookmark",
		api.Handler(
			r.meController.ToggleBookmarkOne(),
			api.HandlerFuncOption{},
		),
	)

	/**
	 * 리뷰 관련
	 */
	rg.POST(
		"/review",
		api.Handler(
			r.meController.WriteReviewOne(),
			api.HandlerFuncOption{},
		),
	)
	rg.PATCH(
		"/review",
		api.Handler(
			r.meController.UpdateRecentReviewOne(),
			api.HandlerFuncOption{},
		),
	)
}

func NewMeRouter(
	meController MeController,
) MeRouter {
	return &meRouter{
		meController: meController,
	}
}
