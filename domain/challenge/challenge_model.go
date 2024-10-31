package challenge

import "time"

type (
	Challenge struct {
		ID          string    `json:"id"`
		UserID      string    `json:"userId"`
		MissionID   string    `json:"missionId"`
		IsCompleted bool      `json:"isCompleted"`
		CompletedAt time.Time `json:"completedAt"`
		CreatedAt   time.Time `json:"createdAt"`
	}
)
