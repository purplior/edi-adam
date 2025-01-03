package logger

import (
	"log"

	"github.com/purplior/podoroot/application/config"
)

func Info(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}

func Error(err error, format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
	if config.DebugMode() {
		log.Println(err.Error())
	}
}

func Debug(format string, v ...interface{}) {
	if config.DebugMode() {
		log.Printf(format+"\n", v...)
	}
}
