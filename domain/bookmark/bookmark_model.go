package bookmark

import (
	"time"

	"github.com/purplior/podoroot/domain/assistant"
)

type (
	Bookmark struct {
		ID          string                  `json:"id"`
		UserID      string                  `json:"userId"`
		AssistantID string                  `json:"assistantId"`
		CreatedAt   time.Time               `json:"createdAt"`
		Assistant   assistant.AssistantInfo `json:"assistantInfo"`
	}
)
