package blip

import (
	"encoding/json"
	"time"
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
	b0 := s[i]
	if b0 < 0x80 {
		e.writeEscapedASCII(buf, b0)
		return 1
	}

	// Manual UTF-8 decode (up to 4 bytes)
	var r rune
	var size int
	switch {
	case b0 < 0xE0 && i+1 < len(s): // 2-byte
		b1 := s[i+1]
		r = rune(b0&0x1F)<<6 | rune(b1&0x3F)
		size = 2
	case b0 < 0xF0 && i+2 < len(s): // 3-byte
		b1, b2 := s[i+1], s[i+2]
		r = rune(b0&0x0F)<<12 | rune(b1&0x3F)<<6 | rune(b2&0x3F)
		size = 3
	case i+3 < len(s): // 4-byte
		b1, b2, b3 := s[i+1], s[i+2], s[i+3]
		r = rune(b0&0x07)<<18 | rune(b1&0x3F)<<12 | rune(b2&0x3F)<<6 | rune(b3&0x3F)
		size = 4
	default:
		// Invalid or truncated
		buf.WriteBytes('\\', 'u', 'f', 'f', 'f', 'd')
		return 1
	}

	buf.WriteRune(r)
	return size
}
