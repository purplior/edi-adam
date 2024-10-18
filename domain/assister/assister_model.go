package assister

import "time"

type (
	AssisterMethod string

	Assister struct {
		ID                 string         `json:"id"`
		ViewID             string         `json:"viewId"`
		AssistantID        string         `json:"assistantId"`
		Method             AssisterMethod `json:"method"`
		AssetURI           string         `json:"assetUri"`
		Version            string         `json:"version"`
		VersionDescription string         `json:"versionDescription"`
		Cost               uint           `json:"cost"`
		CreatedAt          time.Time      `json:"createdAt"`
	}

	// 구현에 대한 정보는 감추고, 보여줘야하는 정보만 보여준다.
	AssisterInfo struct {
		ID                 string `json:"id"`
		Version            string `json:"version"`
		VersionDescription string `json:"versionDescription"`
		CreatedAt          time.Time
	}
)

func (m Assister) ToInfo() AssisterInfo {
	return AssisterInfo{
		ID:                 m.ID,
		Version:            m.Version,
		VersionDescription: m.VersionDescription,
		CreatedAt:          m.CreatedAt,
	}
}

const (
	AssisterMethod_ChatGPT4o AssisterMethod = "gpt4o"
)
