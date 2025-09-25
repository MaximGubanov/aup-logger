package googger

import (
	"context"
	"fmt"
	"googger/pkg/utils"
	"log/slog"
	"os"
)

type CustomLogger struct {
	logger        *slog.Logger
}

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	return m
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m
}

func SetupLogger(logPath, moduleName string, isDebug bool) *CustomLogger {
	logFile, err := utils.CheckingFileExistence(logPath, moduleName)
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %s", err.Error())
	}

	level := slog.LevelInfo
	if isDebug {
		level = slog.LevelDebug
	}

	consoleHandler := NewCustomHandler(os.Stdout, level, moduleName)
	fileHandler := NewCustomHandler(logFile, level, moduleName)

	multiHandler := NewMultiHandler(consoleHandler, fileHandler)

	return &CustomLogger{
		logger: slog.New(multiHandler),
	}
}

func (l *CustomLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *CustomLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *CustomLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *CustomLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *CustomLogger) Fatal(msg string, args ...any) {
	l.logger.Log(nil, LevelFatal, msg, args...)
	os.Exit(1)
}