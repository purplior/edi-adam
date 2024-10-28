package assistant

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	AssistantRepository interface {
		InsertOne(
			ctx inner.Context,
			assistant Assistant,
		) (
			Assistant,
			error,
		)

		FindOneByViewID(
			ctx inner.Context,
			viewID string,
			joinOption AssistantJoinOption,
		) (
			Assistant,
			error,
		)

		FindListByAuthorID(
			ctx inner.Context,
			authorID string,
			joinOption AssistantJoinOption,
		) (
			[]Assistant,
			error,
		)
	}
)
