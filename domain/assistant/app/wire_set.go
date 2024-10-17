package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/assistant"
)

var New = wire.NewSet(
	NewAssistantRouter,
	NewAssistantController,
	domain.NewAssistantService,
)
