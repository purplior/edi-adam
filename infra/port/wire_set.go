package port

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/infra/port/podoopenai"
)

var New = wire.NewSet(
	podoopenai.NewClient,
)
