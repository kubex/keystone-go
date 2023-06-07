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

type EntityLogProvider interface {
	ClearLogs() error
	GetLogs() ([]LogEntry, error)
}

type EntityLogger struct {
	LogEntries []LogEntry `json:",omitempty"`
}

func (e *EntityLogger) ClearLogs() error {
	e.LogEntries = []LogEntry{}
	return nil
}

func (e *EntityLogger) GetLogs() ([]LogEntry, error) {
	return e.LogEntries, nil
}

func (e *EntityLogger) LogDebug(message, reference, traceID string) {
	e.Log(LogLevelDebug, message, reference, traceID, time.Now())
}

func (e *EntityLogger) LogInfo(message, reference, traceID string) {
	e.Log(LogLevelInfo, message, reference, traceID, time.Now())
}

func (e *EntityLogger) LogNotice(message, reference, traceID string) {
	e.Log(LogLevelNotice, message, reference, traceID, time.Now())
}

func (e *EntityLogger) LogWarn(message, reference, traceID string) {
	e.Log(LogLevelWarn, message, reference, traceID, time.Now())
}

func (e *EntityLogger) LogError(message, reference, traceID string) {
	e.Log(LogLevelError, message, reference, traceID, time.Now())
}

func (e *EntityLogger) LogCritical(message, reference, traceID string) {
	e.Log(LogLevelCritical, message, reference, traceID, time.Now())
}

func (e *EntityLogger) LogAlert(message, reference, traceID string) {
	e.Log(LogLevelAlert, message, reference, traceID, time.Now())
}

func (e *EntityLogger) LogFatal(message, reference, traceID string) {
	e.Log(LogLevelFatal, message, reference, traceID, time.Now())
}

func (e *EntityLogger) Log(level LogLevel, message, reference, traceID string, time time.Time) {
	if e.LogEntries == nil {
		e.LogEntries = make([]LogEntry, 0)
	}
	e.LogEntries = append(e.LogEntries, LogEntry{
		written:   false,
		Time:      time,
		Level:     level,
		Message:   message,
		Reference: reference,
		TraceID:   traceID,
	})
}
