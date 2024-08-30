package constants

type (
	Phase int
)

const (
	Phase_Invalid Phase = iota
	Phase_Local
	Phase_Stage
	Phase_Production
)
