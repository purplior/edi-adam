package entity

import (
	"time"

	"github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Assister struct {
		ID                 uint `gorm:"primaryKey;autoIncrement"`
		AssistantID        uint
		Method             assister.AssisterMethod `gorm:"type:varchar(80)"` // 20자 이내
		AssetURI           string                  `gorm:"type:varchar(255)"`
		Version            string                  `gorm:"type:varchar(80)"`
		VersionDescription string                  `gorm:"type:varchar(255)"`
		CreatedAt          time.Time               `gorm:"autoCreateTime"`
	}
)

func (e *Assister) ToModel() assister.Assister {
	model := assister.Assister{
		Method:             e.Method,
		AssetURI:           e.AssetURI,
		Version:            e.Version,
		VersionDescription: e.VersionDescription,
		CreatedAt:          e.CreatedAt,
	}
	if e.ID > 0 {
		model.ID = dt.Str(e.ID)
	}
	if e.AssistantID > 0 {
		model.AssistantID = dt.Str(e.AssistantID)
	}

	return model
}
