package logger_test

import (
	"testing"

	"tagallery.com/api/logger"
)

func TestSetup(t *testing.T) {
	debugLogger := logger.Setup(true)
	prodLogger := logger.Setup(false)

	if debugLogger == logger.Logger() ||
		prodLogger != logger.Logger() {
		t.Error("Logger() should return the last created logger instance.")
	}
}
