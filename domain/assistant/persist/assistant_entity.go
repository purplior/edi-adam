package persist

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Assistant struct {
		ID           uint      `gorm:"primaryKey;autoIncrement"`
		AuthorID     uint      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AuthorID;references:ID"`
		Title        string    `gorm:"type:varchar(80);not null"`  // 20자 이내
		Description  string    `gorm:"type:varchar(255);not null"` // 80자 이내
		VersionLabel string    `gorm:"type:varchar(20);not null"`  // 5자 이내, 공백 미허용
		IsPublic     bool      `gorm:"default:false;not null"`
		CreatedAt    time.Time `gorm:"autoCreateTime"`
	}
)

func (e Assistant) ToModel() domain.Assistant {
	model := domain.Assistant{
		Title:        e.Title,
		Description:  e.Description,
		VersionLabel: e.VersionLabel,
		IsPublic:     e.IsPublic,
		CreatedAt:    e.CreatedAt,
	}

	if e.ID > 0 {
		model.ID = dt.Str(e.ID)
	}
	if e.AuthorID > 0 {
		model.AuthorID = dt.Str(e.AuthorID)
	}

	return model
}

func MakeAssistant(m domain.Assistant) Assistant {
	entity := Assistant{
		Title:        m.Title,
		Description:  m.Description,
		VersionLabel: m.VersionLabel,
		IsPublic:     m.IsPublic,
		CreatedAt:    m.CreatedAt,
	}

	if len(m.ID) > 0 {
		entity.ID = dt.UInt(m.ID)
	}
	if len(m.AuthorID) > 0 {
		entity.AuthorID = dt.UInt(m.AuthorID)
	}

	return entity
}
