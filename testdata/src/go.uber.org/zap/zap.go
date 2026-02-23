package zap

type Logger struct{}

type Field struct{}

func NewProduction() (*Logger, error) { return &Logger{}, nil }

func (l *Logger) Info(msg string, fields ...Field) {}

func (l *Logger) Error(msg string, fields ...Field) {}
