package database

import (
	"log"

	"github.com/purplior/sbec/application/config"
	"github.com/purplior/sbec/domain/shared/exception"
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

	if config.DebugMode() {
		log.Println(err.Error())
	}

	return exception.ErrDBProcess
}
