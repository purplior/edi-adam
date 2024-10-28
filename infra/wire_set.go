package infra

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/port"
	"github.com/podossaem/podoroot/infra/repository"
)

var New = wire.NewSet(
	database.New,
	port.New,
	repository.New,
	NewContextManager,
)
