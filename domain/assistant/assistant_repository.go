package assistant

import "github.com/podossaem/podoroot/domain/context"

type (
	AssistantRepository interface {
		InsertOne(
			ctx context.APIContext,
			assistant Assistant,
		) (
			Assistant,
			error,
		)

		FindListByAuthorID(
			ctx context.APIContext,
			authorID string,
		) (
			[]Assistant,
			error,
		)
	}
)
