package port

import (
	"github.com/google/wire"
	"github.com/purplior/edi-adam/infra/port/openai"
	"github.com/purplior/edi-adam/infra/port/slack"
	"github.com/purplior/edi-adam/infra/port/sms"
)

var New = wire.NewSet(
	openai.NewClient,
	slack.NewClient,
	sms.NewClient,
)
