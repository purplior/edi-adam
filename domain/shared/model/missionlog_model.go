package model

import (
	"time"
)

type (
	MissionLog struct {
		ID         uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		IsAchieved bool `gorm:"default:false;" json:"isAchieved"`
		IsReceived bool `gorm:"default:false;" json:"isReceived"`

		CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
		ReceivedAt *time.Time `json:"receivedAt"`

		UserID    uint     `gorm:"uniqueIndex:idx_user_mission" json:"userId,omitempty"`
		MissionID uint     `gorm:"uniqueIndex:idx_user_mission" json:"missionId,omitempty"`
		Mission   *Mission `json:"mission,omitempty"`
	}
)
