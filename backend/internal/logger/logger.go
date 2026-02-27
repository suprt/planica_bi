package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

// Init initializes the logger based on debug mode
// If isDebug is true: verbose logging with colors (for development)
// If isDebug is false: JSON structured logging, INFO level only (for production)
// If logPath is provided, logs will be written to file with rotation (append mode), otherwise to stdout
func Init(isDebug bool, logPath string, maxBackups int, maxAge int, maxSize int, compress bool) (*zap.Logger, error) {
	var config zap.Config

	if !isDebug {
		// Production mode: JSON format, INFO level only
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		// Development mode: console format with colors, DEBUG level
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// If logPath is provided, write to file with lumberjack rotation
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

		// Use lumberjack for log rotation
		lumberjackLogger := &lumberjack.Logger{
			Filename:   absLogPath,
			MaxSize:    maxSize,    // megabytes
			MaxBackups: maxBackups, // number of old files to keep
			MaxAge:     maxAge,     // days
			Compress:   compress,   // gzip compression
		}

		// Create file writer core
		fileEncoder := config.EncoderConfig
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoder),
			zapcore.AddSync(lumberjackLogger),
			config.Level,
		)

		// Create stdout writer core
		stdoutEncoder := config.EncoderConfig
		if isDebug {
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
