package mission

import (
	"time"

	"github.com/purplior/podoroot/domain/challenge"
)

const (
	MissionReward_Podo3000 MissionReward = "podo_3000"
	MissionReward_Podo5000 MissionReward = "podo_5000"
)

type (
	MissionReward string

	Mission struct {
		ID          string                `json:"id"`
		Title       string                `json:"title"`
		Description string                `json:"description"`
		Reward      MissionReward         `json:"reward"`
		IsPublic    bool                  `json:"isPublic"`
		Challenges  []challenge.Challenge `json:"challenges"`
		CreatedAt   time.Time             `json:"createdAt"`
	}
)
