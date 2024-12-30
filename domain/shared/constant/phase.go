package constant

type (
	Phase int
)

const (
	Phase_Invalid Phase = iota
	Phase_Local
	Phase_Alpha
	Phase_Production
)
