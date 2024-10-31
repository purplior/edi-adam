package mission

import "time"

type (
	MissionReward string

	Mission struct {
		ID                string        `json:"id"`
		Title             string        `json:"string"`
		Description       string        `json:"description"`
		Reward            MissionReward `json:"reward"`
		RewardDescription string        `json:"rewardDescription"`
		IsPublic          bool          `json:"isPublic"`
		CreatedAt         time.Time     `json:"createdAt"`
	}
)

const (
	MissionReward_Podo5000 MissionReward = "podo_5000"
)
