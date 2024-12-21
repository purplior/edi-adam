package entity

import (
	"strings"
	"time"

	domain "github.com/purplior/podoroot/domain/assisterform"
	"github.com/purplior/podoroot/lib/dt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AssisterForm struct {
		ID            primitive.ObjectID     `bson:"_id,omitempty"`
		AssisterID    int                    `bson:"assisterId"`
		Origin        domain.AssisterOrigin  `bson:"origin"`
		Model         domain.AssisterModel   `bson:"model"`
		Fields        []AssisterField        `bson:"fields"`
		Tests         []AssisterInput        `bson:"tests"`
		QueryMessages []AssisterQueryMessage `bson:"queryMessages"`
		CreatedAt     time.Time              `bson:"created_at"`
	}
)

func (e AssisterForm) ToModel() domain.AssisterForm {
	model := domain.AssisterForm{
		Origin:    e.Origin,
		Model:     e.Model,
		CreatedAt: e.CreatedAt,
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

func MakeAssisterForm(m domain.AssisterForm) AssisterForm {
	entity := AssisterForm{
		Origin:    m.Origin,
		Model:     m.Model,
		CreatedAt: m.CreatedAt,
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
		Content []string                        `bson:"content"`
	}
)

func (e AssisterQueryMessage) ToModel() domain.AssisterQueryMessage {
	return domain.AssisterQueryMessage{
		Role:    e.Role,
		Content: strings.Join(e.Content, "\n"),
	}
}

func MakeAssisterQueryMessage(m domain.AssisterQueryMessage) AssisterQueryMessage {
	return AssisterQueryMessage{
		Role:    m.Role,
		Content: strings.Split(m.Content, "\n"),
	}
}

type (
	AssisterInput struct {
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
		Name:   e.Name,
		Values: values,
	}
}

func MakeAssisterInput(m domain.AssisterInput) AssisterInput {
	return AssisterInput{
		Name:   m.Name,
		Values: m.Values,
	}
}
