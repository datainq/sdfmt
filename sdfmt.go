package sdfmt

import (
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/sirupsen/logrus"
	logtypepb "google.golang.org/genproto/googleapis/logging/type"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
)

type StackdriverFormatter struct{}

func (s StackdriverFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	e := logging.Entry{
		Timestamp: entry.Time,
		Payload:   entry.Message,
		Labels:    make(map[string]string),
	}

	switch entry.Level {
	case logrus.PanicLevel:
		e.Severity = logging.Emergency
	case logrus.FatalLevel:
		e.Severity = logging.Emergency
	case logrus.ErrorLevel:
		e.Severity = logging.Critical
	case logrus.WarnLevel:
		e.Severity = logging.Warning
	case logrus.InfoLevel:
		e.Severity = logging.Info
	case logrus.DebugLevel:
		e.Severity = logging.Debug
	case logrus.TraceLevel:
	default:
		e.Severity = logging.Default
	}

	for k, v := range entry.Data {
		e.Labels[k] = fmt.Sprintf("%s", v)
	}

	le, err := toLogEntry(e)
	if err != nil {
		return nil, err
	}

	return json.Marshal(le)
}

func jsonMapToProtoStruct(m map[string]interface{}) *structpb.Struct {
	fields := map[string]*structpb.Value{}
	for k, v := range m {
		fields[k] = jsonValueToStructValue(v)
	}
	return &structpb.Struct{Fields: fields}
}

func jsonValueToStructValue(v interface{}) *structpb.Value {
	switch x := v.(type) {
	case bool:
		return &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: x}}
	case float64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: x}}
	case string:
		return &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: x}}
	case nil:
		return &structpb.Value{Kind: &structpb.Value_NullValue{}}
	case map[string]interface{}:
		return &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: jsonMapToProtoStruct(x)}}
	case []interface{}:
		var vals []*structpb.Value
		for _, e := range x {
			vals = append(vals, jsonValueToStructValue(e))
		}
		return &structpb.Value{Kind: &structpb.Value_ListValue{ListValue: &structpb.ListValue{Values: vals}}}
	default:
		panic(fmt.Sprintf("bad type %T for JSON value", v))
	}
}

func toProtoStruct(v interface{}) (*structpb.Struct, error) {
	// Fast path: if v is already a *structpb.Struct, nothing to do.
	if s, ok := v.(*structpb.Struct); ok {
		return s, nil
	}
	// v is a Go value that supports JSON marshalling. We want a Struct
	// protobuf. Some day we may have a more direct way to get there, but right
	// now the only way is to marshal the Go value to JSON, unmarshal into a
	// map, and then build the Struct proto from the map.
	var jb []byte
	var err error
	if raw, ok := v.(json.RawMessage); ok { // needed for Go 1.7 and below
		jb = []byte(raw)
	} else {
		jb, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("logging: json.Marshal: %v", err)
		}
	}
	var m map[string]interface{}
	err = json.Unmarshal(jb, &m)
	if err != nil {
		return nil, fmt.Errorf("logging: json.Unmarshal: %v", err)
	}
	return jsonMapToProtoStruct(m), nil
}

func toLogEntry(e logging.Entry) (*logpb.LogEntry, error) {
	t := e.Timestamp
	if t.IsZero() {
		t = time.Now().UTC()
	}
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		return nil, err
	}
	ent := &logpb.LogEntry{
		Timestamp:      ts,
		Severity:       logtypepb.LogSeverity(e.Severity),
		InsertId:       e.InsertID,
		Operation:      e.Operation,
		Labels:         e.Labels,
		Trace:          e.Trace,
		SpanId:         e.SpanID,
		Resource:       e.Resource,
		SourceLocation: e.SourceLocation,
	}
	switch p := e.Payload.(type) {
	case string:
		ent.Payload = &logpb.LogEntry_TextPayload{TextPayload: p}
	default:
		s, err := toProtoStruct(p)
		if err != nil {
			return nil, err
		}
		ent.Payload = &logpb.LogEntry_JsonPayload{JsonPayload: s}
	}
	return ent, nil
}
