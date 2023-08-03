package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type EntityLogProvider interface {
	ClearKeystoneLogs() error
	GetKeystoneLogs() ([]*proto.EntityLog, error)
}

type EntityLogger struct {
	ksEntityLogs []*proto.EntityLog
}

func (e *EntityLogger) ClearKeystoneLogs() error {
	e.ksEntityLogs = []*proto.EntityLog{}
	return nil
}

func (e *EntityLogger) GetKeystoneLogs() ([]*proto.EntityLog, error) {
	return e.ksEntityLogs, nil
}

func (e *EntityLogger) LogDebug(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Debug, message, reference, actor, traceID, time.Now(), data)
}

func (e *EntityLogger) LogInfo(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Info, message, reference, actor, traceID, time.Now(), data)
}

func (e *EntityLogger) LogNotice(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Notice, message, reference, actor, traceID, time.Now(), data)
}

func (e *EntityLogger) LogWarn(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Warn, message, reference, actor, traceID, time.Now(), data)
}

func (e *EntityLogger) LogError(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Error, message, reference, actor, traceID, time.Now(), data)
}

func (e *EntityLogger) LogCritical(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Critical, message, reference, actor, traceID, time.Now(), data)
}

func (e *EntityLogger) LogAlert(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Alert, message, reference, actor, traceID, time.Now(), data)
}

func (e *EntityLogger) LogFatal(message, reference, actor, traceID string, data map[string]string) {
	e.KeystoneLog(proto.LogLevel_Fatal, message, reference, actor, traceID, time.Now(), data)
}

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
