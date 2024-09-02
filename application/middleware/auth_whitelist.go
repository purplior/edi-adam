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
		Type   string
		Method string
	}
)

var (
	authWhiteListMap = map[string]AuthWhiteListItem{
		"auth": {
			Type:   Type_AllChildren,
			Method: Method_All,
		},
		"email-verifications": {
			Type:   Type_AllChildren,
			Method: Method_All,
		},
	}
)
