package assister

import (
	"errors"
	"time"
)

const (
	AssisterFieldType_Keyword        AssisterFieldType = "keyword"
	AssisterFieldType_Paragraph      AssisterFieldType = "paragraph"
	AssisterFieldType_ParagraphGroup AssisterFieldType = "paragraph-group"

	QueryMessageRole_System    AssisterQueryMessageRole = "system"
	QueryMessageRole_User      AssisterQueryMessageRole = "user"
	QueryMessageRole_Assistant AssisterQueryMessageRole = "assistant"

	AssisterOrigin_OpenAI               AssisterOrigin = "openai"
	AssisterModel_OpenAI_ChatGPT35Turbo AssisterModel  = "gpt-3.5-turbo"
	AssisterModel_OpenAI_ChatGPT4o      AssisterModel  = "gpt-4o"
	AssisterModel_OpenAI_O1Mini         AssisterModel  = "o1-mini"
)

var (
	ErrInvalidAssisterFieldKey error = errors.New("invalid assister field key")
	ErrInvalidAssisterInput    error = errors.New("invalid assister input")
)

type (
	AssisterFieldType        string
	AssisterQueryMessageRole string
	AssisterOrigin           string
	AssisterModel            string

	AssisterField struct {
		Name   string                 `json:"name"`
		Type   AssisterFieldType      `json:"type"`
		Option map[string]interface{} `json:"option"`
	}

	AssisterRegisterRequest struct {
		Origin        AssisterOrigin         `json:"origin"`
		Model         AssisterModel          `json:"model"`
		Fields        []AssisterField        `json:"fields"`
		Tests         []AssisterInput        `json:"tests"`
		QueryMessages []AssisterQueryMessage `json:"queryMessages"`
		Cost          uint                   `json:"cost"`
	}

	AssisterInput struct {
		Name string `json:"name"`
		/**
		 * Keyword, Paragraph: []string
		 * ParagraphGroup: [][]{
		 *   name string
		 *   value string
		 * }
		 */
		Values []interface{} `json:"values"`
	}
)

func (r AssisterRegisterRequest) ToModalForInsert() Assister {
	return Assister{
		Origin:        r.Origin,
		Model:         r.Model,
		Fields:        r.Fields,
		Tests:         r.Tests,
		QueryMessages: r.QueryMessages,
	}
}

func MakeAssisterFieldFromMap(m map[string]interface{}) (
	AssisterField,
	error,
) {
	field := AssisterField{}
	if name, ok := m["name"].(string); ok {
		field.Name = name
	} else {
		return AssisterField{}, ErrInvalidAssisterFieldKey
	}
	if t, ok := m["type"].(string); ok {
		field.Type = AssisterFieldType(t)
	} else {
		return AssisterField{}, ErrInvalidAssisterFieldKey
	}
	if opt, ok := m["option"].(map[string]interface{}); ok {
		field.Option = opt
	} else {
		return AssisterField{}, ErrInvalidAssisterFieldKey
	}

	return field, nil
}

type (
	Assister struct {
		ID            string                 `json:"id"`
		Origin        AssisterOrigin         `json:"origin"`
		Model         AssisterModel          `json:"model"`
		Fields        []AssisterField        `json:"fields"`
		Tests         []AssisterInput        `json:"tests"`
		QueryMessages []AssisterQueryMessage `json:"queryMessages"`
		Cost          uint                   `json:"cost"`
		CreatedAt     time.Time              `json:"createdAt"`
	}

	AssisterInfo struct {
		ID        string          `json:"id"`
		Fields    []AssisterField `json:"fields"`
		Tests     []AssisterInput `json:"tests"`
		NoStream  bool            `json:"noStream"`
		IsFree    bool            `json:"isFree"`
		Cost      uint            `json:"cost"`
		CreatedAt time.Time       `json:"createdAt"`
	}
)

func (m Assister) Copy() Assister {
	return m
}

func (m Assister) FindField(name string) (AssisterField, bool) {
	for _, field := range m.Fields {
		if field.Name == name {
			return field, true
		}
	}

	return AssisterField{}, false
}

func (m Assister) ToInfo() AssisterInfo {
	noStream := false

	if m.Model == AssisterModel_OpenAI_O1Mini {
		noStream = true
	}

	return AssisterInfo{
		ID:        m.ID,
		Fields:    m.Fields,
		Tests:     m.Tests,
		NoStream:  noStream,
		IsFree:    m.Cost == 0,
		Cost:      m.Cost,
		CreatedAt: m.CreatedAt,
	}
}

type (
	AssisterQueryMessage struct {
		Role    AssisterQueryMessageRole `json:"role"`
		Content string                   `json:"content"`
	}
)
