package assister

type (
	AssisterType string

	Assister struct {
		ID          string       `json:"id"`
		AssistantID string       `json:"assisterId"`
		Type        AssisterType `json:"type"`
		AssetURI    string       `json:"assetUri"`
		Version     string       `json:"version"`
		CreatedAt   string       `json:"createdAt"`
	}
)

const (
	AssisterType_ChatGPT4o AssisterType = "gpt4o"
)
