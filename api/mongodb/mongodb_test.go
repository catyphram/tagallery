package mongodb

import (
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/testutil"
)

func init() {
	config.Load()

	// Manually set the config here for debugging without a config file or env vars
	// config.SetConfig(config.Configuration{
	// 	Database:      "tagallery",
	// 	Database_Host: "localhost:27017",
	// })
}

func TestInitAndClose(t *testing.T) {
	dbOptions := DatabaseOptions{
		Database: config.GetConfig().Database,
		Host:     config.GetConfig().Database_Host,
	}

	db, ctx, err := Init(dbOptions)

	if err != nil {
		format, args := testutil.FormatTestError(
			"Init should return a connection to the mongodb as configured.",
			map[string]interface{}{
				"database": db,
				"context":  ctx,
				"error":    err,
				"options":  dbOptions,
			})
		t.Errorf(format, args...)
	} else {
		Close()
		if _ctx != nil || _db != nil {
			format, args := testutil.FormatTestError(
				"Close should close the db connection and set the variables _ctx and _db to nil.",
				map[string]interface{}{
					"database": _db,
					"context":  _ctx,
				})
			t.Errorf(format, args...)
		}
	}
}

func TestGetConnection(t *testing.T) {
	_, _, err := GetConnection()
	if err == nil {
		format, args := testutil.FormatTestError(
			"GetConnection should return an error when Init() has not been called yet.",
			nil)
		t.Errorf(format, args...)
	}
}
