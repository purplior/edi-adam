package middleware

import "net/http"

const (
	Type_AllChildren AuthWhiteListType = "all_children"
	Type_Exact       AuthWhiteListType = "exact"

	Method_All    AuthWhiteListMethod = "all"
	Method_Get    AuthWhiteListMethod = http.MethodGet
	Method_Post   AuthWhiteListMethod = http.MethodPost
	Method_Put    AuthWhiteListMethod = http.MethodPut
	Method_Delete AuthWhiteListMethod = http.MethodDelete

	Action_Skip         AuthWhiteListAction = "skip"
	Action_SkipButParse AuthWhiteListAction = "skip_but_parse"
	Action_Verify       AuthWhiteListAction = "verify"
)

type (
	AuthWhiteListType   string
	AuthWhiteListMethod string
	AuthWhiteListAction string

	AuthWhiteListItem struct {
		Type     AuthWhiteListType
		Method   AuthWhiteListMethod
		Action   AuthWhiteListAction
		Children map[string]AuthWhiteListItem
	}
)

var (
	authWhiteListMap = map[string]AuthWhiteListItem{
		"auth": {
			Type:   Type_AllChildren,
			Method: Method_All,
			Action: Action_Skip,
		},
		"verifications": {
			Type:   Type_AllChildren,
			Method: Method_All,
			Action: Action_Skip,
		},
		"assistants": {
			Children: map[string]AuthWhiteListItem{
				"detail": {
					Type:   Type_AllChildren,
					Method: Method_Get,
					Action: Action_SkipButParse,
				},
				"podo-list": {
					Type:   Type_AllChildren,
					Method: Method_Get,
					Action: Action_Skip,
				},
			},
		},
		"assisters": {
			Children: map[string]AuthWhiteListItem{
				"exec": {
					Type:   Type_Exact,
					Method: Method_Post,
					Action: Action_SkipButParse,
				},
			},
		},
		"assisterforms": {
			Children: map[string]AuthWhiteListItem{
				"one": {
					Type:   Type_AllChildren,
					Method: Method_Get,
					Action: Action_Skip,
				},
			},
		},
	}
)
