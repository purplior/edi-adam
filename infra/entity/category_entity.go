package entity

import (
	"time"

	domain "github.com/purplior/podoroot/domain/category"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	Category struct {
		ID         uint        `gorm:"primaryKey;autoIncrement"`
		Alias      string      `gorm:"unique;size:30"`
		Label      string      `gorm:"unique;size:50"`
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

	return e
}
