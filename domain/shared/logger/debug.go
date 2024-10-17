package logger

import (
	"log"

	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain/constant"
)

var (
	isDebugLogEnable = config.Phase() != constant.Phase_Production
)

func Debug(format string, v ...interface{}) {
	if isDebugLogEnable {
		log.Printf(format, v...)
	}
}
