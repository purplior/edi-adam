package review

type (
	QueryOption struct {
		ID          uint
		AuthorID    uint
		AssistantID uint
		WithAuthor  bool
	}
)

type (
	AddDTO struct {
		AssistantID uint    `json:"assistantId"`
		Content     string  `json:"content"`
		Score       float64 `json:"score"`
	}
)
