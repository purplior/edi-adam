package port

import (
	"github.com/google/wire"
	"github.com/purplior/podoroot/infra/port/podoopenai"
)

var New = wire.NewSet(
	podoopenai.NewClient,
)
