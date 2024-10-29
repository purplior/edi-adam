package assister

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	AssisterRepository interface {
		FindOneByID(
			ctx inner.Context,
			id string,
		) (
			Assister,
			error,
		)
	}
)
