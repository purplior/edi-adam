package bookmark

import (
	"time"

	"github.com/purplior/sbec/domain/assistant"
)

type (
	Bookmark struct {
		ID            string                  `json:"id"`
		UserID        string                  `json:"userId"`
		AssistantID   string                  `json:"assistantId"`
		CreatedAt     time.Time               `json:"createdAt"`
		AssistantInfo assistant.AssistantInfo `json:"assistantInfo"`
	}
)
