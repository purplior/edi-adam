package entity

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/lib/dt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AssisterForm struct {
		ID               primitive.ObjectID     `bson:"_id,omitempty"`
		AssisterID       int                    `bson:"assisterId"`
		Origin           domain.AssisterOrigin  `bson:"origin"`
		Model            domain.AssisterModel   `bson:"model"`
		Fields           []AssisterField        `bson:"fields"`
		SubmitText       string                 `bson:"submitText"`
		QueryMessages    []AssisterQueryMessage `bson:"queryMessages"`
		QueryInfoHeading string                 `bson:"queryInfoHeading"`
		CreatedAt        time.Time              `bson:"created_at"`
	}
)

func (e AssisterForm) ToModel() domain.AssisterForm {
	model := domain.AssisterForm{
		Origin:           e.Origin,
		Model:            e.Model,
		SubmitText:       e.SubmitText,
		QueryInfoHeading: e.QueryInfoHeading,
		CreatedAt:        e.CreatedAt,
	}

	if !e.ID.IsZero() {
		model.ID = e.ID.Hex()
	}

	if e.AssisterID > 0 {
		model.AssisterID = dt.Str(e.AssisterID)
	}

	model.Fields = make([]domain.AssisterField, len(e.Fields))
	for i, field := range e.Fields {
		model.Fields[i] = field.ToModel()
	}

	model.QueryMessages = make([]domain.AssisterQueryMessage, len(e.QueryMessages))
	for i, queryMessage := range e.QueryMessages {
		model.QueryMessages[i] = queryMessage.ToModel()
	}

	return model
}

func MakeAssisterForm(m domain.AssisterForm) AssisterForm {
	entity := AssisterForm{
		Origin:           m.Origin,
		Model:            m.Model,
		SubmitText:       m.SubmitText,
		QueryInfoHeading: m.QueryInfoHeading,
		CreatedAt:        m.CreatedAt,
	}

	if len(m.ID) > 0 {
		id, _ := primitive.ObjectIDFromHex(m.ID)
		entity.ID = id
	}

	if len(m.AssisterID) > 0 {
		entity.AssisterID = dt.Int(m.AssisterID)
	}

	entity.Fields = make([]AssisterField, len(m.Fields))
	for i, field := range m.Fields {
		entity.Fields[i] = MakeAssisterField(field)
	}

	entity.QueryMessages = make([]AssisterQueryMessage, len(m.QueryMessages))
	for i, queryMessage := range m.QueryMessages {
		entity.QueryMessages[i] = MakeAssisterQueryMessage(queryMessage)
	}

	return entity
}

type (
	AssisterField struct {
		Name     string                   `bson:"name"`
		Type     domain.AssisterFieldType `bson:"type"`
		ItemName string                   `bson:"itemName"`
		Required bool                     `bson:"required"`
		Option   map[string]interface{}   `bson:"option"`
	}
)

func (e AssisterField) ToModel() domain.AssisterField {
	return domain.AssisterField{
		Name:     e.Name,
		Type:     e.Type,
		ItemName: e.ItemName,
		Required: e.Required,
		Option:   e.Option,
	}
}

func MakeAssisterField(m domain.AssisterField) AssisterField {
	return AssisterField{
		Name:     m.Name,
		Type:     m.Type,
		ItemName: m.ItemName,
		Required: m.Required,
		Option:   m.Option,
	}
}

type (
	AssisterQueryMessage struct {
		Role    domain.AssisterQueryMessageRole `bson:"role"`
		Content []string                        `bson:"content"`
	}
)

func (e AssisterQueryMessage) ToModel() domain.AssisterQueryMessage {
	return domain.AssisterQueryMessage{
		Role:    e.Role,
		Content: e.Content,
	}
}

func MakeAssisterQueryMessage(m domain.AssisterQueryMessage) AssisterQueryMessage {
	return AssisterQueryMessage{
		Role:    m.Role,
		Content: m.Content,
	}
}
