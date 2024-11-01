package mission

import "time"

const (
	MissionReward_Podo5000 MissionReward = "podo_5000"
)

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

	MissionInfo struct {
		Title             string        `json:"title"`
		Description       string        `json:"description"`
		Reward            MissionReward `json:"reward"`
		RewardDescription string        `json:"rewardDescription"`
	}
)

func (m Mission) ToInfo() MissionInfo {
	return MissionInfo{
		Title:             m.Title,
		Description:       m.Description,
		Reward:            m.Reward,
		RewardDescription: m.RewardDescription,
	}
}
