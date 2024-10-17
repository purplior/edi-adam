package persist

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/assister"
)

type (
	Assister struct {
		ID                 uint `gorm:"primaryKey;autoIncrement"`
		AssistantID        uint
		Method             domain.AssisterMethod `gorm:"type:varchar(80)"` // 20자 이내
		AssetURI           string                `gorm:"type:varchar(255)"`
		Version            string                `gorm:"type:varchar(80)"`
		VersionDescription string                `gorm:"type:varchar(255)"`
		CreatedAt          time.Time             `gorm:"autoCreateTime"`
	}
)
