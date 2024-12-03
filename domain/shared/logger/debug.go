package logger

import (
	"log"

	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/domain/shared/constant"
)

var (
	isDebugLogEnable = config.Phase() != constant.Phase_Production
)

func Debug(format string, v ...interface{}) {
	if isDebugLogEnable {
		log.Printf(format, v...)
	}
}
