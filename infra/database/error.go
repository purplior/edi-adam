package database

import (
	"log"

	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func ToDomainError(err error) error {
	if err == nil {
		return nil
	}
	if config.DebugMode() {
		log.Println(err.Error())
	}

	switch err {
	case gorm.ErrRecordNotFound:
		return exception.ErrNoRecord
	case mongo.ErrNoDocuments:
		return exception.ErrNoRecord
	}

	return exception.ErrDBProcess
}
