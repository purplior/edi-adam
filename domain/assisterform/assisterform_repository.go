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

		FindOneByID(
			ctx inner.Context,
			id string,
		) (
			AssisterForm,
			error,
		)

		FindOneByAssisterID(
			ctx inner.Context,
			assisterID string,
		) (
			AssisterForm,
			error,
		)
	}
)
