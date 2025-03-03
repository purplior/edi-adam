package router

import "github.com/google/wire"

var New = wire.NewSet(
	NewRouter,
)
