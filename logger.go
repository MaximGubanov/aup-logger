package aup_logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/MaximGubanov/aup-logger/pkg/utils"
)

type CustomLogger struct {
	logger    *slog.Logger
	file      *os.File
	logChan   chan logEntry
	closeChan chan struct{}
	wg        sync.WaitGroup
}

type logEntry struct {
	level string
	msg   string
	args  []any
}

type multiHandler struct {
	handlers []slog.Handler
}

func newMultiHandler(handlers ...slog.Handler) *multiHandler {
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	return m
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m
}

func NewLogger(logPath, moduleName, logLevel string) (*CustomLogger, error) {
	file, err := utils.CheckingFileExistence(logPath, moduleName)
	if err != nil {
		return nil, fmt.Errorf("ошибка настройки логгера: %s", err.Error())
	}

	level := getLogLevel(logLevel)
	consoleHandler := newCustomHandler(os.Stdout, level, moduleName)
	fileHandler := newCustomHandler(file, level, moduleName)
	multiHandler := newMultiHandler(consoleHandler, fileHandler)

	logger := &CustomLogger{
		logger:    slog.New(multiHandler),
		file:      file,
		logChan:   make(chan logEntry),
		closeChan: make(chan struct{}),
	}

	logger.wg.Add(1)
	go logger.processLog()

	return logger, nil
}

func getLogLevel(level string) slog.Level {
	switch level {
	case "D":
		return slog.LevelDebug
	case "E":
		return slog.LevelError
	case "W":
		return slog.LevelWarn
	case "I":
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}

func (l *CustomLogger) processLog() {
	defer l.wg.Done()

	for {
		select {
		case entry := <-l.logChan:
			switch entry.level {
			case "debug":
				l.logger.Debug(entry.msg, entry.args...)
			case "info":
				l.logger.Info(entry.msg, entry.args...)
			case "warn":
				l.logger.Warn(entry.msg, entry.args...)
			case "error":
				l.logger.Error(entry.msg, entry.args...)
			case "fatal":
				l.logger.Error(entry.msg, entry.args...)
				close(l.closeChan)
			}
		case <-l.closeChan:
			for {
				select {
				case entry := <-l.logChan:
					switch entry.level {
					case "info":
						l.logger.Info(entry.msg, entry.args...)
					case "error":
						l.logger.Error(entry.msg, entry.args...)
					}
				default:
					return
				}
			}
		}
	}
}

func (l *CustomLogger) Debug(msg string, args ...any) {
	l.logChan <- logEntry{level: "debug", msg: msg, args: args}
}

func (l *CustomLogger) Info(msg string, args ...any) {
	l.logChan <- logEntry{level: "info", msg: msg, args: args}
}

func (l *CustomLogger) Warn(msg string, args ...any) {
	l.logChan <- logEntry{level: "warn", msg: msg, args: args}
}

func (l *CustomLogger) Error(msg string, args ...any) {
	l.logChan <- logEntry{level: "error", msg: msg, args: args}
}

func (l *CustomLogger) Fatal(msg string, args ...any) {
	l.logChan <- logEntry{level: "fatal", msg: msg, args: args}
}

func (l *CustomLogger) Close() error {
	close(l.closeChan)
	l.wg.Wait()

	return l.file.Close()
}
