package infra

import (
	"github.com/google/wire"
	"github.com/purplior/sbec/infra/database"
	"github.com/purplior/sbec/infra/port"
	"github.com/purplior/sbec/infra/repository"
)

var New = wire.NewSet(
	database.New,
	port.New,
	repository.New,
	NewContextManager,
)
