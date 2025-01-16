package review

import "github.com/purplior/podoroot/domain/shared/inner"

type (
	ReviewRepository interface {
		FindOne_ByAuthorAndAssistantID(
			ctx inner.Context,
			authorID string,
			assistantID string,
			joinOption ReviewJoinOption,
		) (
			Review,
			error,
		)

		InsertOne(
			ctx inner.Context,
			review Review,
		) (
			Review,
			error,
		)

		UpdateOne(
			ctx inner.Context,
			review Review,
		) error
	}
)
