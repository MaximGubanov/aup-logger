package googger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	//"sync"
)

type CustomLogger struct {
	logger *slog.Logger
	file   *os.File
	mu     sync.Mutex
	//logChan   chan logEntry
	//closeChan chan struct{}
	//wg        sync.WaitGroup
}

//type logEntry struct {
//	level string
//	msg   string
//	args  []any
//}

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

func SetupLogger(logPath, moduleName string, isDebug bool) (*CustomLogger, error) {
	file, err := CheckingFileExistence(logPath, moduleName)
	if err != nil {
		return nil, fmt.Errorf("ошибка настройки логгера: %s", err.Error())
	}

	level := slog.LevelInfo
	if isDebug {
		level = slog.LevelDebug
	}

	consoleHandler := NewCustomHandler(os.Stdout, level, moduleName)
	fileHandler := NewCustomHandler(file, level, moduleName)

	multiHandler := NewMultiHandler(consoleHandler, fileHandler)

	logger := &CustomLogger{
		logger: slog.New(multiHandler),
		file:   file,
		mu:     sync.Mutex{},
		//logChan:   make(chan logEntry),
		//closeChan: make(chan struct{}),
	}

	//logger.wg.Add(1)
	//go logger.processLog()

	return logger, nil
}

//func (l *CustomLogger) processLog() {
//	defer l.wg.Done()
//
//	for {
//		select {
//		case entry := <-l.logChan:
//			switch entry.level {
//			case "debug":
//				l.logger.Debug(entry.msg, entry.args...)
//			case "info":
//				l.logger.Info(entry.msg, entry.args...)
//			case "warn":
//				l.logger.Warn(entry.msg, entry.args...)
//			case "error":
//				l.logger.Error(entry.msg, entry.args...)
//			case "fatal":
//				l.logger.Error(entry.msg, entry.args...)
//				close(l.closeChan)
//			}
//		case <-l.closeChan:
//			for {
//				select {
//				case entry := <-l.logChan:
//					switch entry.level {
//					case "info":
//						l.logger.Info(entry.msg, entry.args...)
//					case "error":
//						l.logger.Error(entry.msg, entry.args...)
//					}
//				default:
//					return
//				}
//			}
//		}
//	}
//}

func (l *CustomLogger) Debug(msg string, args ...any) {
	//l.logChan <- logEntry{level: "debug", msg: msg, args: args}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Debug(msg, args...)
}

func (l *CustomLogger) Info(msg string, args ...any) {
	//l.logChan <- logEntry{level: "info", msg: msg, args: args}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Info(msg, args...)
}

func (l *CustomLogger) Warn(msg string, args ...any) {
	//l.logChan <- logEntry{level: "warn", msg: msg, args: args}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Warn(msg, args...)
}

func (l *CustomLogger) Error(msg string, args ...any) {
	//l.logChan <- logEntry{level: "error", msg: msg, args: args}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Error(msg, args...)
}

func (l *CustomLogger) Fatal(msg string, args ...any) {
	//l.logChan <- logEntry{level: "fatal", msg: msg, args: args}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Error(msg, args...)
	os.Exit(1)
}

func (l *CustomLogger) Close() error {
	//close(l.closeChan)
	//l.wg.Wait()

	return l.file.Close()
}
