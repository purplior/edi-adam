package entity

import (
	"time"

	domain "github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/lib/dt"
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

func (e *Assister) ToModel() domain.Assister {
	model := domain.Assister{
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

func MakeAssister(m domain.Assister) Assister {
	e := Assister{
		Version:            m.Version,
		VersionDescription: m.VersionDescription,
		Cost:               m.Cost,
		CreatedAt:          m.CreatedAt,
	}

	if len(m.ID) > 0 {
		e.ID = dt.UInt(m.ID)
	}
	if len(m.AssistantID) > 0 {
		e.AssistantID = dt.UInt(m.AssistantID)
	}

	return e
}
