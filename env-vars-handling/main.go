package main

import (
	"os"

	"github.com/azka-zaydan/article-materials/env-vars-handling/configs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var config *configs.Config

func main() {

	// Set up logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get the config
	config = configs.Get()

	config.Debug()
}
