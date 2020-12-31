package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

// Setup configures the Zap logger.
func Setup(debug bool) *zap.SugaredLogger {
	if debug {
		devLogger, _ := zap.NewDevelopment()
		logger = devLogger.Sugar()
	} else {
		prodLogger, _ := zap.NewProduction()
		logger = prodLogger.Sugar()
	}

	return logger
}

// Logger returns the configured Zap logger. Make sure to call Setup() beforehand.
func Logger() *zap.SugaredLogger {
	return logger
}
