package database

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/infra/database/mymongo"
)

var New = wire.NewSet(
	mymongo.NewClient,
)
