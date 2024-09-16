package database

import (
	"log"

	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain/constant"
	"github.com/podossaem/podoroot/domain/exception"
	"gorm.io/gorm"
)

func ToDomainError(err error) error {
	if config.Phase() == constant.Phase_Local {
		log.Println(err)
	}

	switch err {
	case gorm.ErrRecordNotFound:
		return exception.ErrNoRecord
	}

	return exception.ErrDBProcess
}
