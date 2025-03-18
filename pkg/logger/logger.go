package logger

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init инициализирует глобальный логгер
func Init(serviceName string, level string) {
	once.Do(func() {
		logLevel := getLogLevel(level)

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			logLevel,
		)

		log = zap.New(core,
			zap.AddCaller(),
			zap.Fields(zap.String("service", serviceName)),
		)
	})
}

// getLogLevel преобразует строковый уровень логирования в zapcore.Level
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug логирует сообщение с уровнем Debug
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info логирует сообщение с уровнем Info
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn логирует сообщение с уровнем Warn
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error логирует сообщение с уровнем Error
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal логирует сообщение с уровнем Fatal и завершает программу
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// WithFields возвращает новый Logger с дополнительными полями
func WithFields(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// WithTime добавляет поле времени к логам
func WithTime(name string, t time.Time) zap.Field {
	return zap.Time(name, t)
}

// WithError добавляет ошибку к логам
func WithError(err error) zap.Field {
	return zap.Error(err)
}
