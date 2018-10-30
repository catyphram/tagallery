// Package config loads and provides the configuration options.
// The options may be specified as environment variables or
// a json file can be provided via the command line flag --config.
package config

import (
	"flag"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Database           string `json:"database" env:"DATABASE"`
	Database_Host      string `json:"databaseHost" env:"DATABASE_HOST"`
	Unprocessed_Images string `json:"unprocessedImages" env:"UNPROCESSED_IMAGES"`
	Port               int    `json:"port" env:"PORT"`
}

const (
	configDefault = ""
	configDesc    = "Location of the configuration file"
)

var configuration Configuration
var configPath string
var loaded = false

// parseFlags parses the command line flags and returns the path of the config file
func parseFlags() {
	if !loaded {
		flag.StringVar(&configPath, "c", configDefault, configDesc+" (shorthand)")
		flag.StringVar(&configPath, "config", configDefault, configDesc)
		flag.Parse()
	}

	loaded = true
}

// Load loads the configuration from the config file and env params.
func Load() error {
	parseFlags()
	return gonfig.GetConf(configPath, &configuration)
}

// GetConfig returns the current configuration
func GetConfig() Configuration {
	return configuration
}

// SetConfig sets a new configuration
func SetConfig(config Configuration) {
	configuration = config
}
