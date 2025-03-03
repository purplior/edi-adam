package model

import "time"

type (
	MissionRewardType string
)

var (
	MissionRewardType_Coin MissionRewardType = "coin"
)

type (
	Mission struct {
		ID          uint   `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		Title       string `gorm:"size:255;not null" json:"title"`
		Description string `gorm:"size:400;not null" json:"description"`

		RewardType   MissionRewardType `gorm:"size:20" json:"rewardType"`
		RewardAmount int               `json:"rewardAmount"`
		IsPublished  bool              `json:"isPublished"`

		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

		MissionLogs []MissionLog `gorm:"foreignKey:MissionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"challenges,omitempty"`
	}
)
