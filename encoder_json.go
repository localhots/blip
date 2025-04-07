package blip

import (
	"encoding/json"
	"time"
	"unicode/utf8"
)

type JSONEncoder struct {
	Time       bool
	TimeFormat string
	KeyTime    string
	KeyLevel   string
	KeyMsg     string
}

var _ Encoder = (*JSONEncoder)(nil)

func NewJSONEncoder(cfg Config) JSONEncoder {
	return JSONEncoder{
		Time:       cfg.Time,
		TimeFormat: cfg.TimeFormat,
		KeyTime:    "ts",
		KeyLevel:   "lvl",
		KeyMsg:     "msg",
	}
}

func (e JSONEncoder) EncodeTime(buf *Buffer, ts string) {
	buf.WriteBytes('{')
	e.writeSafeField(buf, e.KeyTime, ts)
	buf.WriteBytes(',')
}

func (e JSONEncoder) EncodeLevel(buf *Buffer, lev Level) {
	if !e.Time {
		buf.WriteBytes('{')
	}
	e.writeSafeField(buf, e.KeyLevel, lev.String())
	buf.WriteBytes(',')
}

func (e JSONEncoder) EncodeMessage(buf *Buffer, msg string) {
	buf.WriteBytes('"')
	buf.WriteString(e.KeyMsg)
	buf.WriteBytes('"', ':', '"')
	buf.WriteString(msg)
	buf.WriteBytes('"')
}

func (e JSONEncoder) EncodeFields(buf *Buffer, lev Level, fields *[]Field) {
	defer buf.WriteBytes('}', '\n')
	if fields == nil || len(*fields) == 0 {
		return
	}

	for _, f := range *fields {
		buf.WriteBytes(',', '"')
		buf.WriteString(f.Key)
		buf.WriteBytes('"', ':')
		e.writeAny(buf, f.Value)
	}
}

func (e JSONEncoder) EncodeStackTrace(buf *Buffer, lev Level, skip int) {
	// Print stack trace but skip the first 4 frames which are part of the
	// logger itself.
	buf.WriteString(stackTrace(skip))
	buf.WriteBytes('\n')
}

// writeSafeField writes a field to the buffer not worrying about escaping it.
func (e JSONEncoder) writeSafeField(buf *Buffer, key, val string) {
	buf.WriteBytes('"')
	buf.WriteString(key)
	buf.WriteBytes('"', ':', '"')
	buf.WriteString(val)
	buf.WriteBytes('"')
}

func (e JSONEncoder) writeAny(buf *Buffer, val any) {
	switch v := val.(type) {
	case string:
		e.writeJSONString(buf, v)
	case []byte:
		panic("TODO: Implement something for []byte")
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
		buf.WriteDuration(v.Truncate(DurationPrecision))
		buf.WriteBytes('"')
	case time.Time:
		buf.WriteBytes('"')
		buf.WriteTime(v, e.TimeFormat)
		buf.WriteBytes('"')
	default:
		if err := json.NewEncoder(buf).Encode(v); err != nil {
			panic(err)
		}
	}
}

// JSON strings

func (e JSONEncoder) writeJSONString(buf *Buffer, s string) {
	buf.WriteBytes('"')

	start := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b < 0x20 || b == '"' || b == '\\' {
			if start < i {
				buf.WriteStringSegment(s, start, i)
			}
			e.writeEscapedASCII(buf, b)
			start = i + 1
		} else if b >= 0x80 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			size := e.writeEscapedUTF8(buf, s, i)
			start = i + size
			i = start - 1
		}
	}

	if start < len(s) {
		buf.WriteStringSegment(s, start, len(s))
	}

	buf.WriteBytes('"')
}

func (e JSONEncoder) writeEscapedASCII(buf *Buffer, b byte) {
	const hex = "0123456789abcdef"
	switch b {
	case '"', '\\':
		buf.WriteBytes('\\', b)
	case '\b':
		buf.WriteBytes('\\', 'b')
	case '\f':
		buf.WriteBytes('\\', 'f')
	case '\n':
		buf.WriteBytes('\\', 'n')
	case '\r':
		buf.WriteBytes('\\', 'r')
	case '\t':
		buf.WriteBytes('\\', 't')
	default:
		buf.WriteBytes('\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
	}
}

func (e JSONEncoder) writeEscapedUTF8(buf *Buffer, s string, i int) int {
	r, size := utf8.DecodeRuneInString(s[i:])
	if r == utf8.RuneError && size == 1 {
		buf.WriteBytes('\\', 'u', 'f', 'f', 'f', 'd')
		return 1
	}
	buf.WriteRune(r)
	return size
}
