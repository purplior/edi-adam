package database

import (
	"github.com/google/wire"
	"github.com/purplior/sbec/infra/database/mongodb"
	"github.com/purplior/sbec/infra/database/redisdb"
	"github.com/purplior/sbec/infra/database/sqldb"
)

var New = wire.NewSet(
	NewDatabaseManager,
	mongodb.NewClient,
	redisdb.NewClient,
	sqldb.NewClient,
)
