package entity

import (
	"time"

	"github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Assister struct {
		ID                 uint                    `gorm:"primaryKey;autoIncrement"`
		AssistantID        uint                    `gorm:"uniqueIndex:idx_assistant_version"`
		AssisterFormID     string                  `gorm:"type:varchar(255)"`
		Method             assister.AssisterMethod `gorm:"type:varchar(80)"` // 20자 이내
		Version            string                  `gorm:"type:varchar(80);uniqueIndex:idx_assistant_version"`
		VersionDescription string                  `gorm:"type:varchar(255)"`
		Cost               uint                    `gorm:"type:tinyint unsigned"`
		CreatedAt          time.Time               `gorm:"autoCreateTime"`
	}
)

func (e *Assister) ToModel() assister.Assister {
	model := assister.Assister{
		Method:             e.Method,
		AssisterFormID:     e.AssisterFormID,
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
