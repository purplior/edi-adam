package assister

import "time"

type (
	Assister struct {
		ID                 string    `json:"id"`
		ViewID             string    `json:"viewId"`
		AssistantID        string    `json:"assistantId"`
		Version            string    `json:"version"`
		VersionDescription string    `json:"versionDescription"`
		Cost               uint      `json:"cost"`
		CreatedAt          time.Time `json:"createdAt"`
	}

	// 구현에 대한 정보는 감추고, 보여줘야하는 정보만 보여준다.
	AssisterInfo struct {
		ID                 string    `json:"id"`
		Version            string    `json:"version"`
		VersionDescription string    `json:"versionDescription"`
		CreatedAt          time.Time `json:"createdAt"`
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
