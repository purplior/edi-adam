package entity

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/mission"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Mission struct {
		ID                uint                 `gorm:"primaryKey;autoIncrement"`
		Title             string               `gorm:"size:255;not null"`
		Description       string               `gorm:"size:400;not null"`
		Reward            domain.MissionReward `gorm:"size:80;not null"`
		RewardDescription string               `gorm:"size:255;not null"`
		IsPublic          bool                 `gorm:"default:false;not null"`
		CreatedAt         time.Time            `gorm:"autoCreateTime"`
		Challenges        []Challenge          `gorm:"foreignKey:MissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (e Mission) ToModel() domain.Mission {
	m := domain.Mission{
		Title:             e.Title,
		Description:       e.Description,
		Reward:            e.Reward,
		RewardDescription: e.RewardDescription,
		IsPublic:          e.IsPublic,
		CreatedAt:         e.CreatedAt,
	}

	if e.ID > 0 {
		m.ID = dt.Str(m.ID)
	}

	return m
}
