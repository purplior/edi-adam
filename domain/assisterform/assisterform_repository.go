package assisterform

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	AssisterFormRepository interface {
		InsertOne(
			ctx inner.Context,
			assisterForm AssisterForm,
		) (
			AssisterForm,
			error,
		)

		FindOne_ByID(
			ctx inner.Context,
			id string,
		) (
			AssisterForm,
			error,
		)

		FindOne_ByAssisterID(
			ctx inner.Context,
			assisterID string,
		) (
			AssisterForm,
			error,
		)
	}
)
