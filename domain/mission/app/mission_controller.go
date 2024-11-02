package app

type (
	MissionController interface{}
)

type (
	missionController struct{}
)

func NewMissionController() MissionController {
	return &missionController{}
}
