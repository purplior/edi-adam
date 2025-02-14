package main

import (
	"github.com/purplior/sbec/application"
	"github.com/purplior/sbec/application/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := application.Start(); err != nil {
		panic(err)
	}
}
