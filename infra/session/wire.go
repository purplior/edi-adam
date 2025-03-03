package session

import (
	"github.com/google/wire"
)

var New = wire.NewSet(
	NewFactory,
)
