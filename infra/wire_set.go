package infra

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/infra/database"
)

var New = wire.NewSet(
	database.New,
)
