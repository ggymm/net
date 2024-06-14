package main

import (
	"net/pkg/app"
	"net/pkg/log"
	"os"
)

func init() {
	app.Init()
	log.Init("hosts")

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("app work dir: %s", wd)
}

func main() {
}
