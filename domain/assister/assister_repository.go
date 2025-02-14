package assister

import "github.com/purplior/sbec/domain/shared/inner"

type (
	AssisterRepository interface {
		FindOne_ByID(
			ctx inner.Context,
			id string,
		) (
			Assister,
			error,
		)

		InsertOne(
			ctx inner.Context,
			assister Assister,
		) (
			Assister,
			error,
		)

		UpdateOne(
			ctx inner.Context,
			assister Assister,
		) error

		DeleteOne_ByID(
			ctx inner.Context,
			id string,
		) error
	}
)
