package main

import (
	"github.com/purplior/podoroot/application"
	"github.com/purplior/podoroot/application/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := application.Start(); err != nil {
		panic(err)
	}
}
