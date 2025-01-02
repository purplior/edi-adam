package category

import "github.com/purplior/podoroot/domain/shared/inner"

type (
	CategoryRepository interface {
		FindOne_ByAlias(
			ctx inner.Context,
			alias string,
		) (
			Category,
			error,
		)

		FindList_ByIDs(
			ctx inner.Context,
			ids []string,
		) (
			[]Category,
			error,
		)
	}
)
