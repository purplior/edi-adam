package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/application/response"
	"github.com/podossaem/podoroot/domain/auth"
	"github.com/podossaem/podoroot/lib/myjwt"
)

var (
	jwtSecretKey = []byte(config.JwtSecretKey())
)

func NewAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			method := request.Method
			pathname := request.URL.Path

			if isWhiteList(method, pathname) {
				return next(c)
			}

			accessToken := request.Header.Get("Authorization")
			accessToken = strings.TrimSpace(strings.Replace(accessToken, "Bearer ", "", 1))

			if len(accessToken) > 0 {
				if identity, err := authorize(accessToken); err == nil {
					return next(&api.Context{
						Context:  c,
						Identity: identity,
					})
				}
			}

			return c.JSON(response.Status_Unauthorized, response.ErrorResponse{
				Status:  response.Status_Unauthorized,
				Message: "unauthorized",
			})
		}
	}
}

func isWhiteList(
	method string,
	pathname string,
) bool {
	// ["", "api", "v1", "auth"]
	segments := strings.Split(pathname, "/")

	if len(segments) < 4 {
		return true
	}

	if item, isItem := authWhiteListMap[segments[3]]; isItem {
		if item.Method == Method_All || item.Method == method {
			return true
		}
	}

	return false
}

func authorize(accessToken string) (*auth.Identity, error) {
	payload, err := myjwt.ParseWithHMAC(accessToken, jwtSecretKey)
	if err != nil {
		return nil, err
	}

	var identity auth.Identity
	identity.SyncWith(payload)

	return &identity, nil
}
