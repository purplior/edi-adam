package assisterform

import "time"

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
	AssisterModel_OpenAI_O1Preview      AssisterModel  = "o1-preview"
)

type (
	AssisterFieldType        string
	AssisterQueryMessageRole string
	AssisterOrigin           string
	AssisterModel            string

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

	AssisterField struct {
		Name     string                 `json:"name"`
		Type     AssisterFieldType      `json:"type"`
		ItemName string                 `json:"itemName"`
		Option   map[string]interface{} `json:"option"`
	}

	AssisterFormRegisterRequest struct {
		Fields        []AssisterField        `json:"fields"`
		QueryMessages []AssisterQueryMessage `json:"queryMessages"`
	}

	AssisterInput struct {
		Type     AssisterFieldType `json:"type"`
		Name     string            `json:"name"`
		ItemName string            `json:"itemName"`
		Values   []interface{}     `json:"values"`
	}
)

type (
	AssisterQueryMessage struct {
		Role    AssisterQueryMessageRole `json:"role"`
		Content string                   `json:"content"`
	}
)

func (m AssisterQueryMessage) CreatePayload() map[string]string {
	return map[string]string{
		"role":    string(m.Role),
		"content": m.Content,
	}
}
