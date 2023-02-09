package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() (*zap.Logger, error) {

	// Get Log variables from ENV
	logLevel := os.Getenv("LOG_LEVEL")
	logFile := os.Getenv("LOG_FILE")

	//fmt.Println(logFile)

	// Check LOG_PATH is not empty
	if logFile == "" {
		logFile = "mylog.log"
	}

	dir := filepath.Dir(logFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating directory: %w", err)
		}
	}

	// Set Config for logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// Set log level
	switch logLevel {
	case "1":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "2":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
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
		return nil, fmt.Errorf("zap error: %w", err)
	}

	return logger, nil
}
