package assistant

import "time"

type (
	Assistant struct {
		ID           string    `json:"id"`
		AuthorID     string    `json:"authorId"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		VersionLabel string    `json:"versionLabel"`
		IsPublic     bool      `json:"isPublic"`
		CreatedAt    time.Time `json:"createdAt"`
	}
)

type (
	RegisterOneRequest struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		VersionLabel string `json:"versionLabel"`
		IsPublic     bool   `json:"isPublic"`
	}
)
