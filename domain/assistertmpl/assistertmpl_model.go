package assistertmpl

type (
	AssisterTemplateField struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Placeholder string `json:"placeholder"`
	}

	AssisterTemplate struct {
		ID        string                  `json:"id"`
		Fields    []AssisterTemplateField `json:"fields"`
		CreatedAt string                  `json:"createdAt"`
	}
)

const (
	AssisterTemplateID_Base = "1"
)
