package log

import (
	"log/slog"
	"os"
)

type Logger struct {
	outLogger *slog.Logger
	errLogger *slog.Logger
}

func NewJsonLogger(level slog.Level) Logger {
	errLogger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
	outLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	return Logger{outLogger: outLogger, errLogger: errLogger}
}

func NewTextLogger(level slog.Level) Logger {
	errLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	outLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	return Logger{outLogger: outLogger, errLogger: errLogger}
}

func (logger *Logger) Debug(msg string, args ...any) {
	logger.outLogger.Debug(msg, args...)
}

func (logger *Logger) Info(msg string, args ...any) {
	logger.outLogger.Info(msg, args...)
}

func (logger *Logger) Warn(msg string, args ...any) {
	logger.errLogger.Warn(msg, args...)
}

func (logger *Logger) Error(msg string, args ...any) {
	logger.errLogger.Error(msg, args...)
}

func (logger *Logger) With(args ...any) Logger {
	return Logger{
		errLogger: logger.errLogger.With(args...),
		outLogger: logger.outLogger.With(args...),
	}
}
