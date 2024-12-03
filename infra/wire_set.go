package infra

import (
	"github.com/google/wire"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/port"
	"github.com/purplior/podoroot/infra/repository"
)

var New = wire.NewSet(
	database.New,
	port.New,
	repository.New,
	NewContextManager,
)
