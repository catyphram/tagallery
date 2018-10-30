package config

import (
	"reflect"
	"testing"

	"tagallery.com/api/testutil"
)

func TestParseFlags(t *testing.T) {
	parseFlags()
	if !loaded {
		format, args := testutil.FormatTestError(
			"Function parseFlags did not set variable loaded to true.",
			map[string]interface{}{
				"loaded": loaded,
			})
		t.Errorf(format, args...)
	}
}

func TestLoad(t *testing.T) {
	Load()
}

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	if !reflect.DeepEqual(configuration, config) {
		format, args := testutil.FormatTestError(
			"Wrong config returned by GetConfig.",
			map[string]interface{}{
				"expected": configuration,
				"got":      config,
			})
		t.Errorf(format, args...)
	}
}

func TestSetConfig(t *testing.T) {
	config := Configuration{
		Database:           "testdb",
		Unprocessed_Images: "testpath",
	}

	SetConfig(config)

	if !reflect.DeepEqual(configuration, config) {
		format, args := testutil.FormatTestError(
			"SetConfig didn't correclty set the config.", map[string]interface{}{
				"expected": config,
				"got":      configuration,
			})
		t.Errorf(format, args...)
	}
}
