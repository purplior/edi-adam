package main

import (
	"github.com/purplior/edi-adam/application"
	"github.com/purplior/edi-adam/application/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := application.Start(); err != nil {
		panic(err)
	}
}
