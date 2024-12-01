package entity

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/category"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Category struct {
		ID         uint   `gorm:"primaryKey;autoIncrement"`
		Alias      string `gorm:"unique;size:30"`
		Label      string `gorm:"size:50"`
		CreatorID  uint
		CreatedAt  time.Time   `gorm:"autoCreateTime"`
		Assistants []Assistant `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (e Category) ToModel() domain.Category {
	m := domain.Category{
		Alias:     e.Alias,
		Label:     e.Label,
		CreatedAt: e.CreatedAt,
	}

	if e.ID > 0 {
		m.ID = dt.Str(e.ID)
	}
	if e.CreatorID > 0 {
		m.CreatorID = dt.Str(e.CreatorID)
	}

	return m
}

func MakeCategory(m domain.Category) Category {
	e := Category{
		Alias:     m.Alias,
		Label:     m.Label,
		CreatedAt: m.CreatedAt,
	}

	if len(m.ID) > 0 {
		e.ID = dt.UInt(m.ID)
	}
	if len(m.CreatorID) > 0 {
		e.CreatorID = dt.UInt(m.CreatorID)
	}

	return e
}
