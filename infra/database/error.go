package database

import (
	"log"
	"runtime/debug"

	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/domain/shared/constant"
	"github.com/purplior/podoroot/domain/shared/exception"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func ToDomainError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return exception.ErrNoRecord
	case mongo.ErrNoDocuments:
		return exception.ErrNoRecord
	}

	if config.Phase() == constant.Phase_Local {
		log.Println(err.Error())
		if err != exception.ErrNoRecord {
			log.Printf("Error: %v\nStack Trace:\n%s", err, debug.Stack())
		}
	}

	return exception.ErrDBProcess
}
