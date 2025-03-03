package controller

import "github.com/purplior/edi-adam/application/common"

type (
	Controller interface {
		// 컨트롤러의 그룹 경로
		GroupPath() string
		// 등록된 핸들러가 자동으로 Router에 바인딩 됨
		Routes() []common.Route
	}
)
