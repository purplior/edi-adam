package model

import (
	"time"

	"github.com/purplior/edi-adam/domain/shared/exception"
)

type (
	AssisterFieldType        string
	AssisterQueryMessageRole string
	AssisterOrigin           string
	AssisterModel            string
)

var (
	AssisterFieldType_Keyword        AssisterFieldType = "keyword"
	AssisterFieldType_Paragraph      AssisterFieldType = "paragraph"
	AssisterFieldType_ParagraphGroup AssisterFieldType = "paragraph-group"

	AssisterQueryMessageRole_System    AssisterQueryMessageRole = "system"
	AssisterQueryMessageRole_User      AssisterQueryMessageRole = "user"
	AssisterQueryMessageRole_Assistant AssisterQueryMessageRole = "assistant"

	AssisterOrigin_OpenAI AssisterOrigin = "openai"
)

type (
	Assister struct {
		ID            string                 `dynamobav:"id" json:"id,omitempty"`
		Origin        AssisterOrigin         `dynamobav:"origin" json:"origin"`
		Model         AssisterModel          `dynamobav:"model" json:"model"`
		Fields        []AssisterField        `dynamobav:"fields" json:"fields"`
		Tests         []AssisterInput        `dynamobav:"tests" json:"tests"`
		QueryMessages []AssisterQueryMessage `dynamobav:"queryMessages" json:"queryMessages"`
		Cost          uint                   `dynamobav:"cost" json:"cost"`

		CreatedAt time.Time `dynamobav:"createdAt" json:"createdAt"`

		AssistantID uint `dynamobav:"assistantId" json:"assistantId"`
	}

	ExcutableAssister struct {
		ID     string          `json:"id"`
		Fields []AssisterField `json:"fields"`
		Tests  []AssisterInput `json:"tests"`
		Cost   uint            `json:"cost"`

		CreatedAt time.Time `json:"createdAt"`
	}
)

func (m Assister) ToExcutable() ExcutableAssister {
	ret := ExcutableAssister{
		ID:     m.ID,
		Fields: m.Fields,
		Tests:  m.Tests,
		Cost:   m.Cost,

		CreatedAt: m.CreatedAt,
	}

	return ret
}

type (
	AssisterField struct {
		Name   string                 `dynamobav:"name" json:"name"`
		Type   AssisterFieldType      `dynamobav:"type" json:"type"`
		Option map[string]interface{} `dynamobav:"option" json:"option"`
	}
)

type (
	AssisterQueryMessage struct {
		Role    AssisterQueryMessageRole `dynamobav:"role" json:"role"`
		Content string                   `dynamobav:"content" json:"content"`
	}
)

type (
	AssisterInput struct {
		Type AssisterFieldType `dynamobav:"type" json:"type"`
		Name string            `dynamobav:"name" json:"name"`
		/**
		 * Keyword, Paragraph: []string
		 * ParagraphGroup: [][]{
		 *   name string
		 *   value string
		 * }
		 */
		Values []interface{} `dynamobav:"values" json:"values"`
	}
)

func MakeAssisterFieldFromMap(m map[string]interface{}) (
	AssisterField,
	error,
) {
	field := AssisterField{}
	if name, ok := m["name"].(string); ok {
		field.Name = name
	} else {
		return AssisterField{}, exception.ErrBadRequest
	}
	if t, ok := m["type"].(string); ok {
		field.Type = AssisterFieldType(t)
	} else {
		return AssisterField{}, exception.ErrBadRequest
	}
	if opt, ok := m["option"].(map[string]interface{}); ok {
		field.Option = opt
	} else {
		return AssisterField{}, exception.ErrBadRequest
	}

	return field, nil
}
