package config_test

import (
	"os"
	"testing"

	"tagallery.com/api/config"
)

func TestLoad(t *testing.T) {
	configuration := config.Load()

	if configuration.Database != "tagallery" {
		t.Error("Load() should set the default settings.")
	}

	if configuration != config.Get() {
		t.Error("Get() should return the same config as Load().")
	}

	os.Setenv("DATABASE_HOST", "database-host")
	os.Setenv("DATABASE", "database")
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "8080")
	os.Setenv("IMAGES", "imgs")

	configuration = config.Load()

	if configuration.DatabaseHost != "database-host" ||
		configuration.Database != "database" ||
		!configuration.Debug ||
		configuration.Port != 8080 ||
		configuration.Images != "imgs" {
		t.Error("Load() should load the settings from the env variables.")
	}
}
