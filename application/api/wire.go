package api

import (
	"github.com/google/wire"
	"github.com/purplior/edi-adam/application/api/controller"
	"github.com/purplior/edi-adam/application/api/router"
)

var New = wire.NewSet(
	controller.New,
	router.New,
)
