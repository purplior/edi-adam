package challenge

import (
	"time"

	"github.com/podossaem/podoroot/domain/mission"
)

type (
	Challenge struct {
		ID         string          `json:"id"`
		UserID     string          `json:"userId"`
		MissionID  string          `json:"missionId"`
		Mission    mission.Mission `json:"mission"`
		IsAchieved bool            `json:"isAchieved"`
		IsReceived bool            `json:"isReceived"`
		ReceivedAt time.Time       `json:"receivedAt"`
		CreatedAt  time.Time       `json:"createdAt"`
	}

	ChallengeInfo struct {
		ID          string              `json:"id"`
		MissionInfo mission.MissionInfo `json:"missionInfo"`
		IsAchieved  bool                `json:"isAchieved"`
		IsReceived  bool                `json:"isReceived"`
		ReceivedAt  time.Time           `json:"receivedAt"`
		CreatedAt   time.Time           `json:"createdAt"`
	}
)

func (m Challenge) ToInfo() ChallengeInfo {
	return ChallengeInfo{
		ID:          m.ID,
		MissionInfo: m.Mission.ToInfo(),
		IsAchieved:  m.IsAchieved,
		IsReceived:  m.IsReceived,
		ReceivedAt:  m.ReceivedAt,
		CreatedAt:   m.CreatedAt,
	}
}
