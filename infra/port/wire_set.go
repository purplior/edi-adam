package port

import (
	"github.com/google/wire"
	"github.com/purplior/podoroot/infra/port/podoopenai"
	"github.com/purplior/podoroot/infra/port/podoslack"
	"github.com/purplior/podoroot/infra/port/podosms"
)

var New = wire.NewSet(
	podoopenai.NewClient,
	podoslack.NewClient,
	podosms.NewClient,
)
