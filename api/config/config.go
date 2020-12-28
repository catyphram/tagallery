package config

import (
	"os"
	"path/filepath"
	"strconv"
)

// Configuration structures all available configuration options.
type Configuration struct {
	Database                string
	DatabaseHost            string
	Debug                   bool
	Port                    int
	Images                  string
	UnprocessedImagesFolder string
	ProcessedImagesFolder   string
}

var config *Configuration

// Get returns the current configuration. Make sure to call Load() beforehand.
func Get() *Configuration {
	return config
}

// Load loads the configuration options from the env variables.
func Load() *Configuration {
	config = &Configuration{
		DatabaseHost:            getEnv("DATABASE_HOST", "localhost:27017"),
		Database:                getEnv("DATABASE", "tagallery"),
		Debug:                   getEnvAsBool("DEBUG", false),
		Port:                    getEnvAsInt("PORT", 3333),
		Images:                  getEnv("IMAGES", filepath.Join(getExecutableDir(), "images")),
		UnprocessedImagesFolder: "unprocessed",
		ProcessedImagesFolder:   "processed",
	}

	return config
}

func getExecutableDir() string {
	if ex, err := os.Executable(); err != nil {
		return filepath.Dir(ex)
	}

	return ""
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(name string, defaultValue bool) bool {
	valStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valStr); err == nil {
		return value
	}
	return defaultValue
}
