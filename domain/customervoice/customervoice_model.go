package customervoice

import "time"

type (
	CustomerVoice struct {
		ID        string    `json:"id"`
		UserID    string    `json:"userId"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

type (
	CustomerVoiceRegisterRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
)
