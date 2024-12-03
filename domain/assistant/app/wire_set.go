package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/podoroot/domain/assistant"
)

var New = wire.NewSet(
	NewAssistantRouter,
	NewAssistantController,
	domain.NewAssistantService,
)
