package sdfmt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type StackdriverFormatter struct{} // TODO(amw): make logEntry fields more configurable.

func (s StackdriverFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	le := logEntry{
		Timestamp:   entry.Time.Format(time.RFC3339),
		Labels:      make(map[string]string),
		TextPayload: entry.Message,
	}

	switch entry.Level {
	case logrus.PanicLevel:
		le.Severity = _alert
	case logrus.FatalLevel:
		le.Severity = _alert
	case logrus.ErrorLevel:
		le.Severity = _critical
	case logrus.WarnLevel:
		le.Severity = _warning
	case logrus.InfoLevel:
		le.Severity = _notice
	case logrus.DebugLevel:
		le.Severity = _info
	case logrus.TraceLevel:
		le.Severity = _debug
	default:
		le.Severity = _default
	}

	for k, v := range entry.Data {
		le.Labels[k] = fmt.Sprintf("%v", v)
	}

	b, err := json.Marshal(le)
	if err != nil {
		return nil, err
	}

	b = append(b, '\n')

	return b, nil
}

type logSeverity int

const (
	_default logSeverity = 100 * iota
	_debug
	_info
	_notice
	_warning
	_error
	_critical
	_alert
	_emergency
)

type logEntry struct {
	LogName      string            `json:"logName,omitempty"`
	Timestamp    string            `json:"timestamp,omitempty"`
	Severity     logSeverity       `json:"severity,omitempty"`
	InsertID     string            `json:"insertId,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Trace        string            `json:"trace,omitempty"`
	SpanID       string            `json:"spanId,omitempty"`
	TraceSampled bool              `json:"traceSampled,omitempty"`
	// TODO(amw): resource, httpRequest, operation, sourceLocation.

	// receiveTimestamp & metadata are output only.

	ProtoPayload interface{}     `json:"protoPayload,omitempty"`
	TextPayload  string          `json:"textPayload,omitempty"`
	JsonPayload  json.RawMessage `json:"jsonPayload,omitempty"`

	// All fields according to
	// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry.
}
