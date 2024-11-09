package assisterform

import (
	"strings"
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

type (
	AssisterFieldType        string
	AssisterQueryMessageRole string
	AssisterOrigin           string
	AssisterModel            string

	AssisterField struct {
		Name     string                 `json:"name"`
		Type     AssisterFieldType      `json:"type"`
		ItemName string                 `json:"itemName"`
		Required bool                   `json:"required"`
		Option   map[string]interface{} `json:"option"`
	}

	AssisterFormRegisterRequest struct {
		Fields        []AssisterField        `json:"fields"`
		QueryMessages []AssisterQueryMessage `json:"queryMessages"`
	}

	AssisterInput struct {
		Name   string        `json:"name"`
		Values []interface{} `json:"values"`
	}
)

type (
	AssisterForm struct {
		ID               string                 `json:"id"`
		AssisterID       string                 `json:"assisterId"`
		Origin           AssisterOrigin         `json:"origin"`
		Model            AssisterModel          `json:"model"`
		Fields           []AssisterField        `json:"fields"`
		SubmitText       string                 `json:"submitText"`
		QueryMessages    []AssisterQueryMessage `json:"queryMessages"`
		QueryInfoHeading string                 `json:"queryInfoHeading"`
		CreatedAt        time.Time              `json:"createdAt"`
	}

	AssisterFormView struct {
		ID         string          `json:"id"`
		AssisterID string          `json:"assisterId"`
		Fields     []AssisterField `json:"fields"`
		SubmitText string          `json:"submitText"`
		NoStream   bool            `json:"noStream"`
		CreatedAt  time.Time       `json:"createdAt"`
	}
)

func (m AssisterForm) FindField(name string) (AssisterField, bool) {
	for _, field := range m.Fields {
		if field.Name == name {
			return field, true
		}
	}

	return AssisterField{}, false
}

func (m AssisterForm) ToView() AssisterFormView {
	noStream := false

	if m.Model == AssisterModel_OpenAI_O1Mini {
		noStream = true
	}

	return AssisterFormView{
		ID:         m.ID,
		AssisterID: m.AssisterID,
		Fields:     m.Fields,
		SubmitText: m.SubmitText,
		NoStream:   noStream,
		CreatedAt:  m.CreatedAt,
	}
}

type (
	AssisterQueryMessage struct {
		Role    AssisterQueryMessageRole `json:"role"`
		Content []string                 `json:"content"`
	}
)

func (m AssisterQueryMessage) CreatePayload() map[string]string {
	return map[string]string{
		"role":    string(m.Role),
		"content": strings.Join(m.Content, "\n"),
	}
}
