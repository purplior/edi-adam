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

		FindOneByID(
			ctx context.APIContext,
			id string,
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
