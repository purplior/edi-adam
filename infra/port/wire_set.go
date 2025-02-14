package port

import (
	"github.com/google/wire"
	"github.com/purplior/sbec/infra/port/openai"
	"github.com/purplior/sbec/infra/port/slack"
	"github.com/purplior/sbec/infra/port/sms"
)

var New = wire.NewSet(
	openai.NewClient,
	slack.NewClient,
	sms.NewClient,
)
