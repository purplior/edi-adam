package assister

import "time"

type (
	AssisterMethod string

	Assister struct {
		ID                 string         `json:"id"`
		AssistantID        string         `json:"assistantId"`
		Method             AssisterMethod `json:"method"`
		AssetURI           string         `json:"assetUri"`
		Version            string         `json:"version"`
		VersionDescription string         `json:"versionDescription"`
		Cost               uint           `json:"cost"`
		CreatedAt          time.Time      `json:"createdAt"`
	}
)

const (
	AssisterMethod_ChatGPT4o AssisterMethod = "gpt4o"
)
