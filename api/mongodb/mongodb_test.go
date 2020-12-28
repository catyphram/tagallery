package mongodb_test

import (
	"context"
	"fmt"
	"testing"

	"tagallery.com/api/config"
	"tagallery.com/api/logger"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/testutil"
)

func init() {
	logger.Setup(true)
}

func TestConnect(t *testing.T) {
	configuration := config.Load()
	client, err := mongodb.Connect(
		context.Background(),
		fmt.Sprintf(`mongodb://%s`, configuration.DatabaseHost),
	)

	if client != mongodb.Client() ||
		err != nil {
		format, args := testutil.FormatTestError(
			"Client() should return a client previously created by Connect().",
			map[string]interface{}{
				"error": err,
			})
		t.Errorf(format, args...)
	}
}
