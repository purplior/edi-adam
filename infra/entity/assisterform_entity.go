package entity

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/lib/dt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AssisterForm struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		AssisterID int                `bson:"assisterId"`
		Fields     []AssisterField    `bson:"fields"`
		SubmitText string             `bson:"submitText"`
		CreatedAt  time.Time          `bson:"created_at"`
	}
)

func (e AssisterForm) ToModel() domain.AssisterForm {
	model := domain.AssisterForm{
		SubmitText: e.SubmitText,
		CreatedAt:  e.CreatedAt,
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

	return model
}

func MakeAssisterForm(m domain.AssisterForm) AssisterForm {
	entity := AssisterForm{
		SubmitText: m.SubmitText,
		CreatedAt:  m.CreatedAt,
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
