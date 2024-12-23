package customervoice

import "github.com/purplior/podoroot/domain/shared/inner"

type (
	CustomerVoiceRepository interface {
		InsertOne(
			ctx inner.Context,
			customerVoice CustomerVoice,
		) (
			CustomerVoice,
			error,
		)
	}
)
