package assisterform

import "time"

var (
	AssisterFieldType_Keyword        = "keyword"
	AssisterFieldType_Paragraph      = "paragraph"
	AssisterFieldType_ParagraphGroup = "paragraph-group"
)

type (
	AssisterFieldType string

	AssisterForm struct {
		ID         string          `json:"id"`
		AssisterID string          `json:"assisterId"`
		Fields     []AssisterField `json:"fields"`
		SubmitText string          `json:"submitText"`
		CreatedAt  time.Time       `json:"createdAt"`
	}

	AssisterField struct {
		Name   string                 `json:"name"`
		Type   AssisterFieldType      `json:"type"`
		Option map[string]interface{} `json:"option"`
	}

	AssisterFormRegisterRequest struct {
		Fields []AssisterField `json:"fields"`
	}
)
