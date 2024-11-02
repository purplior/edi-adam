package assister

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	AssisterRepository interface {
		FindOne_ByID(
			ctx inner.Context,
			id string,
		) (
			Assister,
			error,
		)
	}
)
