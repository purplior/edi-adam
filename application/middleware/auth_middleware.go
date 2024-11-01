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

			action := isWhiteList(method, pathname)
			if action == Action_Skip {
				return next(c)
			}

			accessToken := request.Header.Get("Authorization")
			accessToken = strings.TrimSpace(strings.Replace(accessToken, "Bearer ", "", 1))

			var identity *auth.Identity
			var err error
			if len(accessToken) > 0 {
				identity, err = authorize(accessToken)
			}

			switch action {
			case Action_Verify:
				if err == nil {
					return next(&api.Context{
						Context:  c,
						Identity: identity,
					})
				} else {
					return c.JSON(response.Status_Unauthorized, response.ErrorResponse{
						Status:  response.Status_Unauthorized,
						Message: "unauthorized",
					})
				}
			case Action_SkipAndParse:
				if err == nil {
					return next(&api.Context{
						Context:  c,
						Identity: identity,
					})
				}
			}

			return next(c)
		}
	}
}

func isWhiteList(
	method string,
	pathname string,
) AuthWhiteListAction {
	// ["", "api", "v1", "auth"]
	segments := strings.Split(pathname, "/")

	if len(segments) < 4 {
		return Action_Skip
	}

	if item, isItem := authWhiteListMap[segments[3]]; isItem {
		if item.Method == Method_All || string(item.Method) == method {
			return item.Action
		}
		if len(segments) < 5 {
			return Action_Verify
		}

		children := item.Children
		if len(children) <= 0 {
			return Action_Verify
		}

		if child, isChild := children[segments[4]]; isChild {
			if child.Method == Method_All || string(child.Method) == method {
				return child.Action
			}
		}
	}

	return Action_Verify
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
