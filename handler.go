package aup_logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

const (
	LevelFatal = slog.Level(12)
)

type customHandler struct {
	*slog.Logger
	writer        io.Writer
	level         slog.Level
	nameComponent string
}

func newCustomHandler(w io.Writer, level slog.Level, nc string) *customHandler {
	return &customHandler{writer: w, level: level, nameComponent: nc}
}

func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	timeFormat := r.Time.Format("02.01.2006 15:04:05")
	levelStr := h.formatLevel(r.Level)
	timeParts := strings.Split(timeFormat, " ")
	msg := fmt.Sprintf("[%s][%s][%v][%s][%s]", timeParts[0], timeParts[1], h.nameComponent, levelStr, r.Message)

	if r.NumAttrs() > 0 {
		msg += " |"
		r.Attrs(func(attr slog.Attr) bool {
			msg += fmt.Sprintf(" %s=%v", attr.Key, attr.Value)
			return true
		})
	}

	msg += "\n"
	_, err := h.writer.Write([]byte(msg))

	return err
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *customHandler) formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "D"
	case slog.LevelInfo:
		return "I"
	case slog.LevelWarn:
		return "W"
	case slog.LevelError:
		return "E"
	case LevelFatal:
		return "F"
	default:
		return "UNKNOWN"
	}
}
