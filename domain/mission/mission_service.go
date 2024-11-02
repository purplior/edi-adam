package mission

type (
	MissionService interface{}
)

type (
	missionService struct{}
)

func NewMissionService() MissionService {
	return &missionService{}
}
