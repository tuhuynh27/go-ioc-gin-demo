package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"tuhuynh.com/go-ioc-gin-example/config"
)

type ZapLogger struct {
	Component  struct{}
	Implements struct{}       `implements:"Logger"`
	Config     *config.Config `autowired:"true"`
	logger     *zap.Logger
	sugar      *zap.SugaredLogger
}

func NewZapLogger(config *config.Config) *ZapLogger {
	var logger *zap.Logger
	var err error

	// Use the app mode from the config
	appMode := config.AppMode

	if appMode == "production" {
		// Production config with JSON encoding
		zapConfig := zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err = zapConfig.Build()
	} else {
		// Development config with console encoding
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err = zapConfig.Build()
	}

	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	return &ZapLogger{
		Config: config,
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

// PreDestroy flushes any buffered log entries
func (l *ZapLogger) PreDestroy() {
	if l.logger != nil {
		if err := l.logger.Sync(); err != nil {
			log.Printf("failed to sync logger: %v", err)
		}
	}
}

// Info logs messages at INFO level
func (l *ZapLogger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Debug logs messages at DEBUG level
func (l *ZapLogger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Error logs messages at ERROR level
func (l *ZapLogger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Fatal logs messages at FATAL level and then calls os.Exit(1)
func (l *ZapLogger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// GetLogger returns the initialized Zap logger instance
func (l *ZapLogger) GetLogger() *zap.Logger {
	return l.logger
}

// GetSugar returns the initialized sugared logger instance
func (l *ZapLogger) GetSugar() *zap.SugaredLogger {
	return l.sugar
}

// Sync flushes any buffered log entries
func (l *ZapLogger) Sync() error {
	if l.logger != nil {
		return l.logger.Sync()
	}
	return nil
}
