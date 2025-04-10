package blip

import (
	"fmt"
	"time"
)

// ConsoleEncoder is a console encoder that formats log messages in a
// human-readable format.
type ConsoleEncoder struct {
	Time            bool
	TimeFormat      string
	MinMessageWidth int
	SortFields      bool
	Color           bool
}

const (
	colorRed      = "\033[31m"
	colorGreen    = "\033[32m"
	colorYellow   = "\033[33m"
	colorBlue     = "\033[34m"
	colorPurple   = "\033[35m"
	colorCyan     = "\033[36m"
	colorOffWhite = "\033[37m"
	colorRedBg    = "\033[48;5;88m"
	colorWhite    = "\033[38;5;255m"
	colorReset    = "\033[0m"
)

var _ Encoder = (*ConsoleEncoder)(nil)

// NewConsoleEncoder creates a new console encoder with the given configuration.
// The encoder formats log messages in a human-readable format, with
// colorized levels and optional field sorting.
// The encoder also supports a minimum message width for padding.
func NewConsoleEncoder() ConsoleEncoder {
	return ConsoleEncoder{
		Time:            true,
		TimeFormat:      defaultTimeFormat,
		MinMessageWidth: defaultMessageWidth,
		SortFields:      true,
		Color:           true,
	}
}

// EncodeTime encodes the time of the log message.
func (e ConsoleEncoder) EncodeTime(buf *Buffer, ts string) {
	if !e.Time {
		return
	}
	buf.WriteString(ts)
	buf.WriteBytes(' ')
}

// EncodeLevel encodes the log level of the message.
func (e ConsoleEncoder) EncodeLevel(buf *Buffer, lev Level) {
	e.writeColorized(buf, lev, lev.String())
	buf.WriteBytes(' ')
}

// EncodeMessage encodes the log message.
func (e ConsoleEncoder) EncodeMessage(buf *Buffer, msg string) {
	buf.WriteString(msg)
	if e.MinMessageWidth == 0 {
		return
	}

	// Pad the message to the configured width +2 spaces to separate it from
	// the fields.
	for range e.MinMessageWidth + 2 - len(msg) {
		buf.WriteBytes(' ')
	}
	// If the message is long enough not to be padded, add an extra space to
	// separate it from the fields
	if len(msg) > e.MinMessageWidth {
		// Separate message from fields with 2 spaces
		buf.WriteBytes(' ', ' ')
	}
}

// EncodeFields encodes the fields of the log message.
func (e ConsoleEncoder) EncodeFields(buf *Buffer, lev Level, fields *[]Field) {
	defer buf.WriteBytes('\n')
	if fields == nil || len(*fields) == 0 {
		return
	}
	if e.SortFields {
		sortFields(*fields)
	}

	// Pad fields with two spaces
	buf.WriteBytes(' ', ' ')
	for i, f := range *fields {
		if i > 0 {
			buf.WriteBytes(' ')
		}
		e.writeColorized(buf, lev, f.Key)
		buf.WriteBytes('=')
		e.writeAny(buf, f.Value)
	}
}

// EncodeStackTrace encodes the stack trace of the log message.
func (e ConsoleEncoder) EncodeStackTrace(buf *Buffer, skip int) {
	// Print stack trace but skip the first 4 frames which are part of the
	// logger itself.
	buf.WriteString(stackTrace(skip))
	buf.WriteBytes('\n')
}

// WriteAny writes a value of any type to the buffer. It handles various types
// and falls back to fmt.Sprint for unsupported types.
//
//nolint:gocyclo
func (e ConsoleEncoder) writeAny(buf *Buffer, val any) {
	switch v := val.(type) {
	case string:
		buf.WriteString(v)
	case []byte:
		buf.WriteBytes(v...)
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
	case bool:
		buf.WriteBool(v)
	case time.Duration:
		buf.WriteDuration(v.Truncate(DurationFieldPrecision))
	case time.Time:
		buf.WriteTime(v, TimeFieldFormat)
	default:
		// TODO: Add support for custom encoders
		buf.WriteString(fmt.Sprint(v))
	}
}

func (e ConsoleEncoder) writeColorized(buf *Buffer, lev Level, str string) {
	if !e.Color {
		buf.WriteString(str)
		return
	}

	switch lev {
	case LevelTrace, LevelDebug:
		buf.WriteString(colorOffWhite)
	case LevelInfo:
		buf.WriteString(colorCyan)
	case LevelWarn:
		buf.WriteString(colorYellow)
	case LevelError:
		buf.WriteString(colorRed)
	case LevelPanic, LevelFatal:
		buf.WriteString(colorRedBg)
		buf.WriteString(colorWhite)
	}
	buf.WriteString(str)
	buf.WriteString(colorReset)
}
