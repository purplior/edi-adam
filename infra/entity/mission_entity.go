package entity

import (
	"time"

	"github.com/purplior/podoroot/domain/challenge"
	domain "github.com/purplior/podoroot/domain/mission"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	Mission struct {
		ID          uint                 `gorm:"primaryKey;autoIncrement"`
		Title       string               `gorm:"size:255;not null"`
		Description string               `gorm:"size:400;not null"`
		Reward      domain.MissionReward `gorm:"size:80;not null"`
		IsPublic    bool                 `gorm:"default:false;not null"`
		CreatedAt   time.Time            `gorm:"autoCreateTime"`
		Challenges  []Challenge          `gorm:"foreignKey:MissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (e Mission) ToModel() domain.Mission {
	m := domain.Mission{
		Title:       e.Title,
		Description: e.Description,
		Reward:      e.Reward,
		IsPublic:    e.IsPublic,
		CreatedAt:   e.CreatedAt,
	}

	if e.ID > 0 {
		m.ID = dt.Str(e.ID)
	}

	m.Challenges = make([]challenge.Challenge, len(e.Challenges))
	for i, entity := range e.Challenges {
		m.Challenges[i] = entity.ToModel()
	}

	return m
}

func MakeMission(m domain.Mission) Mission {
	e := Mission{
		Title:       m.Title,
		Description: m.Description,
		Reward:      m.Reward,
		IsPublic:    m.IsPublic,
		CreatedAt:   m.CreatedAt,
	}

	if len(m.ID) > 0 {
		e.ID = dt.UInt(m.ID)
	}

	return e
}
