package blip

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

// Buffer is a byte buffer used for encoding log entries.
// WARNING: Buffer should not be initialized manually. It is pooled to reduce
// allocations.
type Buffer struct {
	b []byte
}

const bufferSize = 1024

//
// Encoding
//

func (buf *Buffer) WriteAny(val any) {
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
		buf.WriteDuration(v.Truncate(DurationPrecision))
	case time.Time:
		buf.WriteTime(v, TimeFormat)
	default:
		// TODO: Add support for custom encoders
		buf.WriteString(fmt.Sprint(v))
	}
}

// Write implements the io.Writer interface.
func (buf *Buffer) Write(b []byte) (int, error) {
	buf.b = append(buf.b, b...)
	return len(b), nil
}

// WriteBytes writes a byte slice to the buffer.
func (buf *Buffer) WriteBytes(b ...byte) {
	buf.b = append(buf.b, b...)
}

// WriteString writes a string to the buffer.
func (buf *Buffer) WriteString(str string) {
	buf.b = append(buf.b, str...)
}

func (buf *Buffer) WriteRune(r rune) {
	var tmp [utf8.UTFMax]byte
	n := utf8.EncodeRune(tmp[:], r)
	buf.b = append(buf.b, tmp[:n]...)
}

// WriteInt writes an int64 value to the buffer.
func (buf *Buffer) WriteInt(i int64) {
	buf.b = strconv.AppendInt(buf.b, i, 10)
}

// WriteUint writes a uint64 value to the buffer.
func (buf *Buffer) WriteUint(i uint64) {
	buf.b = strconv.AppendUint(buf.b, i, 10)
}

// WriteFloat writes a float64 value to the buffer with the specified bit size.
func (buf *Buffer) WriteFloat(f float64, bitSize int) {
	buf.b = strconv.AppendFloat(buf.b, f, 'f', -1, bitSize)
}

// WriteBool writes a boolean value to the buffer.
func (buf *Buffer) WriteBool(b bool) {
	buf.b = strconv.AppendBool(buf.b, b)
}

// WriteDuration writes a time.Duration value to the buffer.
func (buf *Buffer) WriteDuration(d time.Duration) {
	buf.b = append(buf.b, d.String()...)
}

// WriteTime writes a time.Time value to the buffer using the specified format.
func (buf *Buffer) WriteTime(t time.Time, format string) {
	buf.b = t.AppendFormat(buf.b, format)
}

func (buf *Buffer) WriteEscapedString(s string) {
	buf.WriteBytes('"')

	start := 0
	for i := 0; i < len(s); {
		b := s[i]

		// ASCII needing escape
		if b < 0x20 || b == '"' || b == '\\' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.writeEscapedASCII(b)
			i++
			start = i
			continue
		}

		// Probably not ASCII
		if b >= 0x80 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			size := buf.writeEscapedUTF8(s, i)
			i += size
			start = i
			continue
		}

		// No escape needed
		i++
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}

	buf.WriteBytes('"')
}

func (buf *Buffer) writeEscapedASCII(b byte) {
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

func (buf *Buffer) writeEscapedUTF8(s string, i int) int {
	r, size := utf8.DecodeRuneInString(s[i:])
	if r == utf8.RuneError && size == 1 {
		buf.WriteBytes('\\', 'u', 'f', 'f', 'f', 'd')
		return 1
	}
	buf.WriteRune(r)
	return size
}

func (buf *Buffer) WriteBase64(data []byte) {
	buf.WriteBytes('"')

	// Calculate the needed size for base64 encoding
	needed := base64.StdEncoding.EncodedLen(len(data))

	// Check if the buffer has enough capacity
	if cap(buf.b)-len(buf.b) < needed {
		// Not enough space, append the required amount
		buf.b = append(buf.b, make([]byte, needed)...)
	}

	// Encode the data into the buffer (now safely large enough)
	start := len(buf.b)
	base64.StdEncoding.Encode(buf.b[start:], data)

	buf.WriteBytes('"')
}

//
// Buffer pool
//

// Buffers are pooled to reduce allocations.
var bufferPool = sync.Pool{
	New: func() any {
		buf := Buffer{make([]byte, 0, bufferSize)}
		return &buf
	},
}

func getBuffer() *Buffer {
	s, _ := bufferPool.Get().(*Buffer)
	s.b = s.b[:0] // Reset the underlying slice
	return s
}

func putBuffer(buf *Buffer) {
	bufferPool.Put(buf)
}
