package blip

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// JSONEncoder is an encoder that encodes log messages in JSON format.
type JSONEncoder struct {
	TimeFormat     string
	TimePrecision  time.Duration
	Base64Encoding *base64.Encoding
	KeyTime        string
	KeyLevel       string
	KeyMsg         string
	KeyStackTrace  string

	timeCache func(time.Time) string
}

var _ Encoder = (*JSONEncoder)(nil)

// NewJSONEncoder creates a new JSON encoder with the given configuration.
// The encoder formats log messages in JSON format, with optional and fields.
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		TimeFormat:     defaultTimeFormat,
		TimePrecision:  defaultTimePrecision,
		Base64Encoding: base64.StdEncoding,
		KeyTime:        "ts",
		KeyLevel:       "lvl",
		KeyMsg:         "msg",
		KeyStackTrace:  "stack",
	}
}

// Start writes the beginning of the log message.
func (e *JSONEncoder) Start(buf *Buffer) {
	buf.WriteBytes('{')
}

// EncodeTime encodes the time of the log message.
func (e *JSONEncoder) EncodeTime(buf *Buffer) {
	if e.TimeFormat == "" {
		return
	}
	if e.timeCache == nil && e.TimePrecision > 0 {
		e.timeCache = timeCache(e.TimeFormat, e.TimePrecision)
	}
	if e.timeCache != nil {
		e.writeSafeField(buf, e.KeyTime, e.timeCache(timeNow()))
	} else {
		buf.WriteBytes('"')
		buf.WriteTime(timeNow(), e.TimeFormat)
		buf.WriteBytes('"')
	}
}

// EncodeLevel encodes the log level of the message.
func (e *JSONEncoder) EncodeLevel(buf *Buffer, lev Level) {
	if e.TimeFormat != "" {
		buf.WriteBytes(',')
	}
	e.writeSafeField(buf, e.KeyLevel, lev.String())
}

// EncodeMessage encodes the log message.
func (e *JSONEncoder) EncodeMessage(buf *Buffer, msg string) {
	buf.WriteBytes(',', '"')
	buf.WriteString(e.KeyMsg)
	buf.WriteBytes('"', ':')
	buf.WriteEscapedString(msg)
}

// EncodeFields encodes the fields of the log message.
func (e *JSONEncoder) EncodeFields(buf *Buffer, _ Level, fields *[]Field) {
	if fields == nil || len(*fields) == 0 {
		return
	}

	for _, f := range *fields {
		buf.WriteBytes(',')
		buf.WriteEscapedString(f.Key)
		buf.WriteBytes(':')
		e.writeAny(buf, f.Value)
	}
}

// EncodeStackTrace encodes the stack trace of the log message.
func (e *JSONEncoder) EncodeStackTrace(buf *Buffer, skip int) {
	buf.WriteBytes(',', '"')
	buf.WriteString(e.KeyStackTrace)
	buf.WriteBytes('"', ':')
	buf.WriteEscapedString(stackTrace(skip))
}

// End writes the end of the log message.
func (e *JSONEncoder) End(buf *Buffer) {
	buf.WriteBytes('}', '\n')
}

// writeSafeField writes a field to the buffer not worrying about escaping it.
func (e *JSONEncoder) writeSafeField(buf *Buffer, key, val string) {
	buf.WriteBytes('"')
	buf.WriteString(key)
	buf.WriteBytes('"', ':', '"')
	buf.WriteString(val)
	buf.WriteBytes('"')
}

// nolint:gocyclo
func (e *JSONEncoder) writeAny(buf *Buffer, val any) {
	switch v := val.(type) {
	case string:
		buf.WriteEscapedString(v)
	case []byte:
		buf.WriteBase64(e.Base64Encoding, v)
	case nil:
		buf.WriteString("null")
	case bool:
		buf.WriteBool(v)
	case int:
		buf.WriteInt(int64(v))
	case int8:
		buf.WriteInt(int64(v))
	case int16:
		buf.WriteInt(int64(v))
	case int32:
		buf.WriteInt(int64(v))
	case int64:
		buf.WriteInt(v)
	case uint:
		buf.WriteUint(uint64(v))
	case uint8:
		buf.WriteUint(uint64(v))
	case uint16:
		buf.WriteUint(uint64(v))
	case uint32:
		buf.WriteUint(uint64(v))
	case uint64:
		buf.WriteUint(v)
	case float32:
		buf.WriteFloat(float64(v), 32)
	case float64:
		buf.WriteFloat(v, 64)
	case time.Duration:
		buf.WriteBytes('"')
		buf.WriteDuration(v.Truncate(DurationFieldPrecision))
		buf.WriteBytes('"')
	case time.Time:
		buf.WriteBytes('"')
		buf.WriteTime(v, TimeFieldFormat)
		buf.WriteBytes('"')
	default:
		//nolint:errchkjson
		_ = json.NewEncoder(buf).Encode(v)
	}
}
