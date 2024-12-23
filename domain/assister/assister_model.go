package assister

import (
	"errors"
	"time"
)

var (
	ErrInvalidAssisterInput error = errors.New("invalid assister input")
)

type (
	Assister struct {
		ID                 string    `json:"id"`
		AssistantID        string    `json:"assistantId"`
		Version            string    `json:"version"`
		VersionDescription string    `json:"versionDescription"`
		Cost               uint      `json:"cost"`
		CreatedAt          time.Time `json:"createdAt"`
	}

	// 구현에 대한 정보는 감추고, 보여줘야하는 정보만 보여준다.
	AssisterInfo struct {
		ID                 string    `json:"id"`
		IsFree             bool      `json:"isFree"`
		Version            string    `json:"version"`
		VersionDescription string    `json:"versionDescription"`
		Cost               uint      `json:"cost"`
		CreatedAt          time.Time `json:"createdAt"`
	}
)

func (m Assister) ToInfo() AssisterInfo {
	return AssisterInfo{
		ID:                 m.ID,
		IsFree:             m.Cost == 0,
		Version:            m.Version,
		VersionDescription: m.VersionDescription,
		Cost:               m.Cost,
		CreatedAt:          m.CreatedAt,
	}
}

type (
	AssisterRegisterRequest struct {
		AssistantID        string `json:"assistantId"`
		Version            string `json:"version"`
		VersionDescription string `json:"versionDescription"`
		Cost               uint   `json:"cost"`
	}
)
