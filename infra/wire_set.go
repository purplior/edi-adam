package infra

import (
	"github.com/google/wire"
	"github.com/podossaem/root/infra/database"
)

var New = wire.NewSet(
	database.New,
)
