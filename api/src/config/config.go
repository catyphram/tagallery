package config

import (
	"flag"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Unprocessed_Images string
}

var configuration Configuration

func LoadConfig() error {
	var configPath string

	const (
		configDefault = "config.json"
		configDesc    = "Location of the configuration file"
	)
	flag.StringVar(&configPath, "c", configDefault, configDesc+" (shorthand)")
	flag.StringVar(&configPath, "config", configDefault, configDesc)
	flag.Parse()

	err := gonfig.GetConf(configPath, &configuration)
	return err
}

func GetConfig() Configuration {
	return configuration
}
