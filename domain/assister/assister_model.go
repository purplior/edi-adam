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

	Origin_OpenAI = "openai"
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
		Type AssisterFieldType `json:"type"`
		Name string            `json:"name"`
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

func (r AssisterRegisterRequest) ToModelForInsert() Assister {
	return Assister{
		Origin:        r.Origin,
		Model:         r.Model,
		Fields:        r.Fields,
		Tests:         r.Tests,
		QueryMessages: r.QueryMessages,
		Temperature:   0.7,
		TopP:          0.8,
		Cost:          r.Cost,
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
		Temperature   float64                `json:"temperature"`
		TopP          float64                `json:"topP"`
		Cost          uint                   `json:"cost"`
		CreatedAt     time.Time              `json:"createdAt"`
	}

	AssisterInfo struct {
		ID        string          `json:"id"`
		Fields    []AssisterField `json:"fields"`
		Tests     []AssisterInput `json:"tests"`
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
	return AssisterInfo{
		ID:        m.ID,
		Fields:    m.Fields,
		Tests:     m.Tests,
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
