package entity

import (
	"time"

	"github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Assister struct {
		ID                 uint      `gorm:"primaryKey;autoIncrement"`
		AssistantID        uint      `gorm:"uniqueIndex:idx_assistant_version"`
		Version            string    `gorm:"size:80;uniqueIndex:idx_assistant_version"`
		VersionDescription string    `gorm:"size:255"`
		Cost               uint      `gorm:"type:tinyint unsigned"`
		CreatedAt          time.Time `gorm:"autoCreateTime"`
	}
)

func (e *Assister) ToModel() assister.Assister {
	model := assister.Assister{
		Version:            e.Version,
		VersionDescription: e.VersionDescription,
		Cost:               e.Cost,
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
