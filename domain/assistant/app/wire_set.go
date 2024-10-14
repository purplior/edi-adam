package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/assistant/persist"
)

var New = wire.NewSet(
	NewAssistantRouter,
	NewAssistantController,
	domain.NewAssistantService,
	persist.NewAssistantRepository,
)
