package entity

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/assisterform"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AssisterForm struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		AssisterID uint               `json:"assisterId"`
		Fields     []AssisterField    `bson:"_fields"`
		CreatedAt  time.Time          `bson:"created_at"`
	}
)

func (e AssisterForm) ToModel() domain.AssisterForm {
	id := ""
	if !e.ID.IsZero() {
		id = e.ID.Hex()
	}

	fields := make([]domain.AssisterField, len(e.Fields))
	for i, field := range e.Fields {
		fields[i] = field.ToModel()
	}

	return domain.AssisterForm{
		ID:        id,
		Fields:    fields,
		CreatedAt: e.CreatedAt,
	}
}

func MakeAssisterForm(m domain.AssisterForm) AssisterForm {
	var id primitive.ObjectID
	if len(m.ID) == 0 {
		id, _ = primitive.ObjectIDFromHex(m.ID)
	}

	fields := make([]AssisterField, len(m.Fields))
	for i, field := range m.Fields {
		fields[i] = MakeAssisterField(field)
	}

	return AssisterForm{
		ID:        id,
		Fields:    fields,
		CreatedAt: m.CreatedAt,
	}
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
