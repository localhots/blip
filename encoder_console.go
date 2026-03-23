package blip

import (
	"fmt"
	"strings"
	"time"
)

// ConsoleEncoder is a console encoder that formats log messages in a
// human-readable format.
type ConsoleEncoder struct {
	TimeFormat      string
	TimePrecision   time.Duration
	MinMessageWidth int
	SortFields      bool
	Color           bool

	timeCache   func(time.Time) string
	levelLabels [7]string // Pre-built colorized level labels
	levelColors [7]string // Color prefix per level
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
	fontBold      = "\033[1m"
	fontReset     = "\033[0m"
)

var _ Encoder = (*ConsoleEncoder)(nil)

var (
	padding    = strings.Repeat(" ", 128)
	levelNames = [7]string{"TRAC", "DEBU", "INFO", "WARN", "ERRO", "PANI", "FATA"}
)

// NewConsoleEncoder creates a new console encoder with the given configuration.
// The encoder formats log messages in a human-readable format, with
// colorized levels and optional field sorting.
// The encoder also supports a minimum message width for padding.
func NewConsoleEncoder() *ConsoleEncoder {
	return &ConsoleEncoder{
		TimeFormat:      defaultTimeFormat,
		TimePrecision:   defaultTimePrecision,
		MinMessageWidth: defaultMessageWidth,
		SortFields:      true,
		Color:           true,
	}
}

// prepare pre-computes cached values. Called from Logger.New().
func (e *ConsoleEncoder) prepare() {
	// Pre-build level labels (with trailing space) and color prefixes
	colors := [7]string{
		colorOffWhite,           // Trace
		colorOffWhite,           // Debug
		colorCyan,               // Info
		colorYellow,             // Warn
		colorRed,                // Error
		colorRedBg + colorWhite, // Panic
		colorRedBg + colorWhite, // Fatal
	}
	for i, name := range levelNames {
		e.levelColors[i] = colors[i]
		if e.Color {
			e.levelLabels[i] = colors[i] + name + fontReset + " "
		} else {
			e.levelLabels[i] = name + " "
		}
	}

	// Initialize time cache (includes trailing space)
	if e.TimeFormat != "" && e.TimePrecision > 0 {
		e.timeCache = timeCache("", " ", e.TimeFormat, e.TimePrecision)
	}
}

// Start writes the beginning of the log message.
func (e *ConsoleEncoder) Start(_ *Buffer) {}

// EncodeTime encodes the time of the log message.
func (e *ConsoleEncoder) EncodeTime(buf *Buffer) {
	if e.timeCache != nil {
		buf.WriteString(e.timeCache(timeNow()))
		return
	}
	e.encodeTimeSlow(buf)
}

func (e *ConsoleEncoder) encodeTimeSlow(buf *Buffer) {
	if e.TimeFormat == "" {
		return
	}
	buf.WriteTime(timeNow(), e.TimeFormat)
	buf.WriteBytes(' ')
}

// EncodeLevel encodes the log level of the message.
func (e *ConsoleEncoder) EncodeLevel(buf *Buffer, lev Level) {
	buf.WriteString(e.levelLabels[lev-1])
}

// EncodeMessage encodes the log message.
func (e *ConsoleEncoder) EncodeMessage(buf *Buffer, msg string) {
	if e.Color {
		buf.WriteString(fontBold)
	}
	buf.WriteString(msg)
	if e.Color {
		buf.WriteString(fontReset)
	}
	if e.MinMessageWidth > 0 {
		e.padMessage(buf, msg)
	}
}

func (e *ConsoleEncoder) padMessage(buf *Buffer, msg string) {
	// Pad the message to the configured width +2 spaces to separate it from
	// the fields.
	if padLen := e.MinMessageWidth + 2 - len(msg); padLen > 0 {
		buf.WriteString(padding[:padLen])
	}
	// If the message is long enough not to be padded, add an extra space to
	// separate it from the fields
	if len(msg) > e.MinMessageWidth {
		// Separate message from fields with 2 spaces
		buf.WriteBytes(' ', ' ')
	}
}

// EncodeFields encodes the fields of the log message.
func (e *ConsoleEncoder) EncodeFields(buf *Buffer, lev Level, fields *[]Field) {
	if fields == nil || len(*fields) == 0 {
		return
	}
	if e.SortFields {
		sortFields(*fields)
	}

	color := e.levelColors[lev-1]
	useColor := e.Color

	// Pad fields with two spaces
	buf.WriteBytes(' ', ' ')
	for i, f := range *fields {
		if i > 0 {
			buf.WriteBytes(' ')
		}
		if useColor {
			buf.WriteString(color)
			buf.WriteString(f.Key)
			buf.WriteString(fontReset)
		} else {
			buf.WriteString(f.Key)
		}
		buf.WriteBytes('=')
		e.writeAny(buf, f.Value)
	}
}

// EncodeStackTrace encodes the stack trace of the log message.
func (e *ConsoleEncoder) EncodeStackTrace(buf *Buffer, skip int) {
	buf.WriteBytes('\n')
	writeStackTrace(buf, skip)
}

// End writes the end of the log message.
func (e *ConsoleEncoder) End(buf *Buffer) {
	buf.WriteBytes('\n')
}

// WriteAny writes a value of any type to the buffer. It handles various types
// and falls back to fmt.Sprint for unsupported types.
//
//nolint:gocyclo
func (e *ConsoleEncoder) writeAny(buf *Buffer, val any) {
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
	case error:
		buf.WriteString(v.Error())
	case nil:
		buf.WriteString("<nil>")
	default:
		buf.WriteString(fmt.Sprint(v))
	}
}
