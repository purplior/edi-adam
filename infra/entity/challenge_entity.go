package entity

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/challenge"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Challenge struct {
		ID         uint `gorm:"primaryKey;autoIncrement"`
		UserID     uint
		MissionID  uint
		IsAchieved bool `gorm:"default:false;not null"`
		IsReceived bool `gorm:"default:false;not null"`
		ReceivedAt time.Time
		CreatedAt  time.Time `gorm:"autoCreateTime"`
	}
)

func (e Challenge) ToModel() domain.Challenge {
	m := domain.Challenge{
		IsAchieved: e.IsAchieved,
		IsReceived: e.IsReceived,
		ReceivedAt: e.ReceivedAt,
		CreatedAt:  e.CreatedAt,
	}

	if e.ID > 0 {
		m.ID = dt.Str(e.ID)
	}
	if e.UserID > 0 {
		m.UserID = dt.Str(e.UserID)
	}
	if e.MissionID > 0 {
		m.MissionID = dt.Str(e.MissionID)
	}

	return m
}

func MakeChallenge(m domain.Challenge) Challenge {
	e := Challenge{
		IsAchieved: m.IsAchieved,
		IsReceived: m.IsReceived,
		ReceivedAt: m.ReceivedAt,
		CreatedAt:  m.CreatedAt,
	}

	if len(m.ID) > 0 {
		e.ID = dt.UInt(m.ID)
	}
	if len(m.UserID) > 0 {
		e.UserID = dt.UInt(m.UserID)
	}
	if len(m.MissionID) > 0 {
		e.MissionID = dt.UInt(m.MissionID)
	}

	return e
}
