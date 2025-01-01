package category

import "github.com/purplior/podoroot/domain/shared/inner"

type (
	CategoryRepository interface {
		FindList_ByIDs(
			ctx inner.Context,
			ids []string,
		) (
			[]Category,
			error,
		)
	}
)
