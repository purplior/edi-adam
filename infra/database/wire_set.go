package database

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/infra/database/mymongo"
	"github.com/podossaem/podoroot/infra/database/myredis"
)

var New = wire.NewSet(
	mymongo.NewClient,
	myredis.NewClient,
)
