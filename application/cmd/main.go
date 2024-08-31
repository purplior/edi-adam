package main

import (
	"github.com/podossaem/podoroot/application"
	"github.com/podossaem/podoroot/application/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := application.Start(); err != nil {
		panic(err)
	}
}
