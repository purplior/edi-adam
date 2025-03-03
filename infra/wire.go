package infra

import (
	"github.com/google/wire"
	"github.com/purplior/edi-adam/infra/database"
	"github.com/purplior/edi-adam/infra/port"
	"github.com/purplior/edi-adam/infra/repository"
	"github.com/purplior/edi-adam/infra/session"
)

var New = wire.NewSet(
	database.New,
	port.New,
	repository.New,
	session.New,
)
