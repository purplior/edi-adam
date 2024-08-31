package domain

import (
	"github.com/google/wire"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

var New = wire.NewSet(
	verification.New,
)
