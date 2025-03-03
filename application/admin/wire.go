package admin

import (
	"github.com/google/wire"
	"github.com/purplior/edi-adam/application/admin/controller"
	"github.com/purplior/edi-adam/application/admin/router"
)

var New = wire.NewSet(
	controller.New,
	router.New,
)
