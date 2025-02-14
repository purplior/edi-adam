package entity

import (
	"time"

	domain "github.com/purplior/sbec/domain/assister"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Assister struct {
		ID            primitive.ObjectID     `bson:"_id,omitempty"`
		Origin        string                 `bson:"origin"`
		Model         string                 `bson:"model"`
		Fields        []AssisterField        `bson:"fields"`
		Tests         []AssisterInput        `bson:"tests"`
		QueryMessages []AssisterQueryMessage `bson:"queryMessages"`
		Cost          uint                   `bson:"cost"`
		Temperature   float64                `bson:"temperature"`
		TopP          float64                `bson:"topP"`
		CreatedAt     time.Time              `bson:"created_at"`
	}
)

func (e Assister) ToModel() domain.Assister {
	model := domain.Assister{
		Origin:      domain.AssisterOrigin(e.Origin),
		Model:       domain.AssisterModel(e.Model),
		Cost:        e.Cost,
		Temperature: e.Temperature,
		TopP:        e.TopP,
		CreatedAt:   e.CreatedAt,
	}

	if !e.ID.IsZero() {
		model.ID = e.ID.Hex()
	}

	model.Fields = make([]domain.AssisterField, len(e.Fields))
	for i, field := range e.Fields {
		model.Fields[i] = field.ToModel()
	}

	model.Tests = make([]domain.AssisterInput, len(e.Tests))
	for i, test := range e.Tests {
		model.Tests[i] = test.ToModel()
	}

	model.QueryMessages = make([]domain.AssisterQueryMessage, len(e.QueryMessages))
	for i, queryMessage := range e.QueryMessages {
		model.QueryMessages[i] = queryMessage.ToModel()
	}

	return model
}

func MakeAssister(m domain.Assister) Assister {
	entity := Assister{
		Origin:      string(m.Origin),
		Model:       string(m.Model),
		Cost:        m.Cost,
		Temperature: m.Temperature,
		TopP:        m.TopP,
		CreatedAt:   m.CreatedAt,
	}

	if len(m.ID) > 0 {
		id, _ := primitive.ObjectIDFromHex(m.ID)
		entity.ID = id
	}

	entity.Fields = make([]AssisterField, len(m.Fields))
	for i, field := range m.Fields {
		entity.Fields[i] = MakeAssisterField(field)
	}

	entity.Tests = make([]AssisterInput, len(m.Tests))
	for i, test := range m.Tests {
		entity.Tests[i] = MakeAssisterInput(test)
	}

	entity.QueryMessages = make([]AssisterQueryMessage, len(m.QueryMessages))
	for i, queryMessage := range m.QueryMessages {
		entity.QueryMessages[i] = MakeAssisterQueryMessage(queryMessage)
	}

	return entity
}

type (
	AssisterField struct {
		Name   string                   `bson:"name"`
		Type   domain.AssisterFieldType `bson:"type"`
		Option map[string]interface{}   `bson:"option"`
	}
)

func (e AssisterField) ToModel() domain.AssisterField {
	return domain.AssisterField{
		Name:   e.Name,
		Type:   e.Type,
		Option: e.Option,
	}
}

func MakeAssisterField(m domain.AssisterField) AssisterField {
	return AssisterField{
		Name:   m.Name,
		Type:   m.Type,
		Option: m.Option,
	}
}

type (
	AssisterQueryMessage struct {
		Role    domain.AssisterQueryMessageRole `bson:"role"`
		Content string                          `bson:"content"`
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

type (
	AssisterInput struct {
		Type   string        `bson:"type"`
		Name   string        `bson:"name"`
		Values []interface{} `bson:"values"`
	}
)

func (e AssisterInput) ToModel() domain.AssisterInput {
	values := make([]interface{}, len(e.Values))
	for i, value := range e.Values {
		switch v := value.(type) {
		case string:
			values[i] = v
		case bson.A:
			vGroups := make([]map[string]string, len(v))
			for j, vItem := range v {
				bItem := vItem.(bson.D)

				item := map[string]string{
					bItem[0].Key: bItem[0].Value.(string),
					bItem[1].Key: bItem[1].Value.(string),
				}

				vGroups[j] = item
			}
			values[i] = vGroups
		}
	}

	return domain.AssisterInput{
		Type:   domain.AssisterFieldType(e.Type),
		Name:   e.Name,
		Values: values,
	}
}

func MakeAssisterInput(m domain.AssisterInput) AssisterInput {
	return AssisterInput{
		Type:   string(m.Type),
		Name:   m.Name,
		Values: m.Values,
	}
}
