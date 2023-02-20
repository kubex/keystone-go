package keystone

import (
	"time"
)

type LogEntry struct {
	written   bool
	Time      time.Time
	Level     LogLevel
	Message   string // Log message
	Reference string // File Name/Line
	TraceID   string // Trace ID / Correlation ID
	Actor     *Actor
}

func (e *Entity) Debug(message, reference, traceID string) {
	e.Log(LogLevelDebug, message, reference, traceID, time.Now())
}

func (e *Entity) Info(message, reference, traceID string) {
	e.Log(LogLevelInfo, message, reference, traceID, time.Now())
}

func (e *Entity) Notice(message, reference, traceID string) {
	e.Log(LogLevelNotice, message, reference, traceID, time.Now())
}

func (e *Entity) Warn(message, reference, traceID string) {
	e.Log(LogLevelWarn, message, reference, traceID, time.Now())
}

func (e *Entity) Error(message, reference, traceID string) {
	e.Log(LogLevelError, message, reference, traceID, time.Now())
}

func (e *Entity) Critical(message, reference, traceID string) {
	e.Log(LogLevelCritical, message, reference, traceID, time.Now())
}

func (e *Entity) Alert(message, reference, traceID string) {
	e.Log(LogLevelAlert, message, reference, traceID, time.Now())
}

func (e *Entity) Fatal(message, reference, traceID string) {
	e.Log(LogLevelFatal, message, reference, traceID, time.Now())
}

func (e *Entity) Log(level LogLevel, message, reference, traceID string, time time.Time) {
	e.LogEntries = append(e.LogEntries, LogEntry{
		written:   false,
		Time:      time,
		Level:     level,
		Message:   message,
		Reference: reference,
		TraceID:   traceID,
	})
}
