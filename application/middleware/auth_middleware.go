package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/purplior/edi-adam/application/common"
	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/lib/myjwt"
)

var (
	jwtSecretKey = []byte(config.JwtSecretKey())
)

func NewAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(_ctx echo.Context) error {
			c := _ctx.(*common.Context)
			request := c.Request()

			accessToken := request.Header.Get("Authorization")
			accessToken = strings.TrimSpace(strings.Replace(accessToken, "Bearer ", "", 1))

			var identity *inner.Identity = nil
			var err error = nil
			if len(accessToken) > 0 {
				identity, err = authorize(accessToken)
				if err == nil && identity != nil {
					c.Identity = identity
				}
			}

			return next(c)
		}
	}
}

func authorize(accessToken string) (*inner.Identity, error) {
	payload, err := myjwt.ParseWithHMAC(accessToken, jwtSecretKey)
	if err != nil {
		return nil, err
	}

	var identity inner.Identity
	identity.SyncWith(payload)

	return &identity, nil
}
