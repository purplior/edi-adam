package database

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/infra/database/mymongo"
	"github.com/podossaem/podoroot/infra/database/myredis"
	"github.com/podossaem/podoroot/infra/database/podosql"
)

var New = wire.NewSet(
	NewDatabaseManager,
	mymongo.NewClient,
	myredis.NewClient,
	podosql.NewClient,
)
