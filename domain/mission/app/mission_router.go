package app

type (
	MissionRouter interface{}
)

type (
	missionRouter struct{}
)

func NewMissionRouter() MissionRouter {
	return &missionRouter{}
}
