package middleware

import "net/http"

const (
	Type_AllChildren = "allChildren"
	Type_Exact       = "exact"

	Method_All    = "all"
	Method_Get    = http.MethodGet
	Method_Post   = http.MethodPost
	Method_Put    = http.MethodPut
	Method_Delete = http.MethodDelete
)

type (
	AuthWhiteListItem struct {
		Type     string
		Method   string
		Children map[string]AuthWhiteListItem
	}
)

var (
	authWhiteListMap = map[string]AuthWhiteListItem{
		"auth": {
			Type:   Type_AllChildren,
			Method: Method_All,
		},
		"verifications": {
			Type:   Type_AllChildren,
			Method: Method_All,
		},
		"assistants": {
			Children: map[string]AuthWhiteListItem{
				"detail": {
					Type:   Type_AllChildren,
					Method: Method_Get,
				},
				"podo-list": {
					Type:   Type_AllChildren,
					Method: Method_Get,
				},
			},
		},
		"assisters": {
			Children: map[string]AuthWhiteListItem{
				"exec": {
					Type:   Type_Exact,
					Method: Method_Post,
				},
			},
		},
		"assisterforms": {
			Type:   Type_AllChildren,
			Method: Method_Get,
		},
	}
)
