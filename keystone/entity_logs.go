package keystone

import (
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EntityLogProvider is an interface for entities that can have logs
type EntityLogProvider interface {
	ClearKeystoneLogs() error
	GetKeystoneLogs() []*proto.EntityLog
}

// EntityLogger is a struct that implements EntityLogProvider
type EntityLogger struct {
	ksEntityLogs []*proto.EntityLog
}

// ClearKeystoneLogs clears the logs
func (e *EntityLogger) ClearKeystoneLogs() error {
	e.ksEntityLogs = []*proto.EntityLog{}
	return nil
}

// GetKeystoneLogs returns the logs
func (e *EntityLogger) GetKeystoneLogs() []*proto.EntityLog {
	return e.ksEntityLogs
}

// LogDebug logs a debug message
func (e *EntityLogger) LogDebug(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Debug, message, reference, actor, traceID, time.Now(), data)
}

// LogInfo logs an info message
func (e *EntityLogger) LogInfo(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Info, message, reference, actor, traceID, time.Now(), data)
}

// LogNotice logs a notice message
func (e *EntityLogger) LogNotice(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Notice, message, reference, actor, traceID, time.Now(), data)
}

// LogWarn logs a warning message
func (e *EntityLogger) LogWarn(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Warn, message, reference, actor, traceID, time.Now(), data)
}

// LogError logs an error message
func (e *EntityLogger) LogError(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Error, message, reference, actor, traceID, time.Now(), data)
}

// LogCritical logs a critical message
func (e *EntityLogger) LogCritical(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Critical, message, reference, actor, traceID, time.Now(), data)
}

// LogAlert logs an alert message
func (e *EntityLogger) LogAlert(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Alert, message, reference, actor, traceID, time.Now(), data)
}

// LogFatal logs a fatal message
func (e *EntityLogger) LogFatal(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Fatal, message, reference, actor, traceID, time.Now(), data)
}

// KeystoneLog logs a message
func (e *EntityLogger) KeystoneLog(level proto.LogLevel, message, reference, actor, traceID string, logTime time.Time, data map[string]string) {
	if e.ksEntityLogs == nil {
		e.ksEntityLogs = make([]*proto.EntityLog, 0)
	}
	e.ksEntityLogs = append(e.ksEntityLogs, &proto.EntityLog{
		Actor:     actor,
		Level:     level,
		Message:   message,
		Reference: reference,
		TraceId:   traceID,
		Time:      timestamppb.New(logTime),
		Data:      data,
	})
}
