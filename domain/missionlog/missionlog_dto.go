package missionlog

type (
	QueryOption struct {
		ID          uint
		UserID      uint
		MissionID   uint
		WithMission bool
	}

	AchieveDTO struct {
		UserID    uint
		MissionID uint
	}
)
