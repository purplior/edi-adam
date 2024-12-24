package category

import "github.com/purplior/podoroot/domain/shared/inner"

type (
	CategoryRepository interface {
		FindListByIDs(
			ctx inner.Context,
			ids []string,
		) (
			[]Category,
			error,
		)
	}
)
