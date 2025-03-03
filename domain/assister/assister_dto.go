package assister

import (
	"errors"

	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/lib/validator"
)

var (
	ErrInvalidAssisterFieldKey error = errors.New("invalid assister field key")
	ErrInvalidAssisterInput    error = errors.New("invalid assister input")
)

type (
	QueryOption struct {
		ID string
	}

	RegisterDTO struct {
		AssistantID   uint
		ID            string
		Origin        model.AssisterOrigin
		Model         model.AssisterModel
		Fields        []model.AssisterField
		QueryMessages []model.AssisterQueryMessage
		Tests         []model.AssisterInput
	}

	UpdateDTO struct {
		ID string

		Origin        model.AssisterOrigin
		Model         model.AssisterModel
		Fields        []model.AssisterField
		QueryMessages []model.AssisterQueryMessage
		Tests         []model.AssisterInput
	}

	RequestDTO struct {
		ID     string
		UserID uint
		Inputs []model.AssisterInput
	}
)

func (m RegisterDTO) IsValid() bool {
	if !validator.CheckValidAssisterFields(m.Fields) {
		return false
	}
	if !validator.CheckValidAssisterQueryMessages(
		m.Fields,
		m.QueryMessages,
	) {
		return false
	}

	return true
}
