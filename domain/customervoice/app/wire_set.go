package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/customervoice"
)

var New = wire.NewSet(
	NewCustomerVoiceController,
	NewCustomerVoiceRouter,
	domain.NewCustomerVoiceService,
)
