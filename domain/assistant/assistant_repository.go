package assistant

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	AssistantRepository interface {
		InsertOne(
			ctx context.APIContext,
			assistant Assistant,
		) (
			Assistant,
			error,
		)

		FindOneByViewID(
			ctx context.APIContext,
			viewID string,
			joinOption AssistantJoinOption,
		) (
			Assistant,
			error,
		)

		FindListByAuthorID(
			ctx context.APIContext,
			authorID string,
			joinOption AssistantJoinOption,
		) (
			[]Assistant,
			error,
		)
	}
)
