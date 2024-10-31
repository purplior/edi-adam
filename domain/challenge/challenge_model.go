package challenge

import (
	"time"

	"github.com/podossaem/podoroot/domain/mission"
)

type (
	Challenge struct {
		ID          string          `json:"id"`
		UserID      string          `json:"userId"`
		MissionID   string          `json:"missionId"`
		Mission     mission.Mission `json:"mission"`
		IsCompleted bool            `json:"isCompleted"`
		CompletedAt time.Time       `json:"completedAt"`
		CreatedAt   time.Time       `json:"createdAt"`
	}
)
