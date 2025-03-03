package database

import (
	"github.com/google/wire"
	"github.com/purplior/edi-adam/infra/database/dynamo"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

var New = wire.NewSet(
	NewDatabaseManager,
	dynamo.NewClient,
	postgre.NewClient,
)
