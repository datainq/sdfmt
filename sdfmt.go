package sdfmt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type StackdriverFormatter struct{}

func (s StackdriverFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	le := logEntry{
		Labels:    make(map[string]string),
		Message:   entry.Message,
		Timestamp: entry.Time.Format(time.RFC3339Nano),
	}

	switch entry.Level {
	case logrus.PanicLevel:
		le.Severity = _emergency
	case logrus.FatalLevel:
		le.Severity = _emergency
	case logrus.ErrorLevel:
		le.Severity = _error
	case logrus.WarnLevel:
		le.Severity = _warning
	case logrus.InfoLevel:
		le.Severity = _info
	case logrus.DebugLevel:
		le.Severity = _debug
	case logrus.TraceLevel:
		le.Severity = _debug
	default:
		le.Severity = _alert
	}

	for k, v := range entry.Data {
		le.Labels[k] = fmt.Sprintf("%v", v)
	}

	data, err := json.Marshal(le)
	return append(data, '\n'), err
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
	Labels    map[string]string `json:"labels,omitempty"`
	Message   string            `json:"message,omitempty"`
	Severity  logSeverity       `json:"severity,omitempty"`
	Timestamp string            `json:"timestamp,omitempty"`
}
