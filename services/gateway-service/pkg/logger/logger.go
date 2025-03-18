package logger

// Logger интерфейс для логирования
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

// Заглушка для логирования
type mockLogger struct{}

// NewMockLogger создает новый экземпляр заглушки логгера
func NewMockLogger() Logger {
	return &mockLogger{}
}

func (l *mockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *mockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *mockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *mockLogger) Error(msg string, keysAndValues ...interface{}) {}
func (l *mockLogger) Fatal(msg string, keysAndValues ...interface{}) {}
