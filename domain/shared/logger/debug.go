package logger

import (
	"log"

	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/domain/shared/constant"
)

var (
	isDebugLogEnable = config.Phase() != constant.Phase_Production
)

func Info(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}

func Error(err error, format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
	if isDebugLogEnable {
		log.Println(err.Error())
	}
}

func Debug(format string, v ...interface{}) {
	if isDebugLogEnable {
		log.Printf(format+"\n", v...)
	}
}
