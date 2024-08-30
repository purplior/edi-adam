package main

import (
	"github.com/podossaem/root/application"
	"github.com/podossaem/root/application/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := application.Start(); err != nil {
		panic(err)
	}
}
