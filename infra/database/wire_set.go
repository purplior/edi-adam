package database

import (
	"github.com/google/wire"
	"github.com/purplior/podoroot/infra/database/myredis"
	"github.com/purplior/podoroot/infra/database/podomongo"
	"github.com/purplior/podoroot/infra/database/podosql"
)

var New = wire.NewSet(
	NewDatabaseManager,
	podomongo.NewClient,
	myredis.NewClient,
	podosql.NewClient,
)
