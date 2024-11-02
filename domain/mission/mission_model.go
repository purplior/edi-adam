package mission

import (
	"time"

	"github.com/podossaem/podoroot/domain/challenge"
)

const (
	MissionReward_Podo5000 MissionReward = "podo_5000"
)

type (
	MissionReward string

	Mission struct {
		ID                string                `json:"id"`
		Title             string                `json:"string"`
		Description       string                `json:"description"`
		Reward            MissionReward         `json:"reward"`
		RewardDescription string                `json:"rewardDescription"`
		IsPublic          bool                  `json:"isPublic"`
		Challenges        []challenge.Challenge `json:"challenges"`
		CreatedAt         time.Time             `json:"createdAt"`
	}
)
