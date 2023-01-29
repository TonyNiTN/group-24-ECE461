package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() (*zap.Logger, error) {

	// Get Log variables from ENV
	logLevel := os.Getenv("LOG_LEVEL")
	logFile := os.Getenv("LOG_FILE")

	// Check LOG_PATH is not empty
	if logFile == "" {
		logFile = "mylog.log"
	}

	// Set Config for logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Set log level
	switch logLevel {
	case "1":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "2":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel + 1)
	}

	// Set output file
	if logFile != "" {
		config.OutputPaths = []string{logFile}
		config.ErrorOutputPaths = []string{logFile}
	}

	// Create logger
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
