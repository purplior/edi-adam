package assister

import "time"

type (
	AssisterMethod string

	Assister struct {
		ID                 string         `json:"id"`
		Method             AssisterMethod `json:"method"`
		AssetURI           string         `json:"assetUri"`
		Version            string         `json:"version"`
		VersionDescription string         `json:"versionDescription"`
		CreatedAt          time.Time      `json:"createdAt"`
		AssistantID        string         `json:"assistantId"`
	}
)

const (
	AssisterMethod_ChatGPT4o AssisterMethod = "gpt4o"
)
