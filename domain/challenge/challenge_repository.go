package challenge

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	ChallengeRepository interface {
		FindPaginatedListByUserID(
			ctx inner.Context,
			userId string,
			limit int,
			offset int,
		) (
			[]Challenge,
			error,
		)
	}
)
