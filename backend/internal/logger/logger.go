package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init initializes the logger based on environment
// If logPath is provided, logs will be written to file (append mode), otherwise to stdout
func Init(env, logPath string) (*zap.Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// If logPath is provided, write to file
	if logPath != "" {
		// Resolve path to absolute to avoid issues with relative paths
		absLogPath, err := filepath.Abs(logPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve log path: %w", err)
		}

		// Create directory if it doesn't exist
		dir := filepath.Dir(absLogPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open file in append mode (O_APPEND) to prevent overwriting on restart
		// If file doesn't exist, it will be created
		logFile, err := os.OpenFile(absLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}

		// Create file writer core
		fileEncoder := config.EncoderConfig
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoder),
			zapcore.AddSync(logFile),
			config.Level,
		)

		// Create stdout writer core
		stdoutEncoder := config.EncoderConfig
		if env != "production" {
			stdoutEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		stdoutCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(stdoutEncoder),
			zapcore.AddSync(os.Stdout),
			config.Level,
		)

		// Combine cores: write to both file and stdout
		core := zapcore.NewTee(fileCore, stdoutCore)

		logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		Log = logger
		return logger, nil
	}

	// No file path - write only to stdout
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	Log = logger
	return logger, nil
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
