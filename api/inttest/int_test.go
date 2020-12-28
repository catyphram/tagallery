package inttest

import (
	"testing"
	"time"

	"tagallery.com/api/server"
)

// startAPI boots to API.
func startAPI() {
	go server.StartServer()

	// Wait 100 milliseconds to ensure that the http server has started
	// before the integration tests are run. A race condition may otherwise occur.
	time.Sleep(100 * time.Millisecond)
}

func TestAPI(t *testing.T) {
	startAPI()

	t.Run("GetCategories", GetCategories)
	t.Run("GetImages", GetImages)
	t.Run("PostCategory", PostCategory)
	t.Run("PostImage", PostImage)
}
