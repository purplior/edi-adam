package challenge

import (
	"time"
)

type (
	Challenge struct {
		ID         string    `json:"id"`
		UserID     string    `json:"userId"`
		MissionID  string    `json:"missionId"`
		IsAchieved bool      `json:"isAchieved"`
		IsReceived bool      `json:"isReceived"`
		ReceivedAt time.Time `json:"receivedAt"`
		CreatedAt  time.Time `json:"createdAt"`
	}
)
