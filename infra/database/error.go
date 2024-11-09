package database

import (
	"log"
	"runtime/debug"

	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain/shared/constant"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"gorm.io/gorm"
)

func ToDomainError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return exception.ErrNoRecord
	}

	if config.Phase() == constant.Phase_Local {
		log.Println(err)
		if err != exception.ErrNoRecord {
			log.Printf("Error: %v\nStack Trace:\n%s", err, debug.Stack())
		}
	}

	return exception.ErrDBProcess
}

func ToReadError(err error) string {
	if err == nil {
		return ReadError_Success
	}

	switch err {
	case gorm.ErrRecordNotFound:
		return ReadError_NoRecord
	}

	return ReadError_ErrorOccurred
}

var (
	ReadError_NoRecord      = "no_record"
	ReadError_ErrorOccurred = "error_occurred"
	ReadError_Success       = "success"
)
