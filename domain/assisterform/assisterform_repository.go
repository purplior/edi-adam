package assisterform

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	AssisterFormRepository interface {
		InsertOne(
			ctx context.APIContext,
			assisterForm AssisterForm,
		) (
			AssisterForm,
			error,
		)

		FindOneByID(
			ctx context.APIContext,
			id string,
		) (
			AssisterForm,
			error,
		)
	}
)
