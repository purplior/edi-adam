package database

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/infra/database/myredis"
	"github.com/podossaem/podoroot/infra/database/podomongo"
	"github.com/podossaem/podoroot/infra/database/podopaysql"
	"github.com/podossaem/podoroot/infra/database/podosql"
)

var New = wire.NewSet(
	NewDatabaseManager,
	podomongo.NewClient,
	myredis.NewClient,
	podosql.NewClient,
	podopaysql.NewClient,
)
