package blip

import (
	"encoding/base64"
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

// WriteEscapedString writes a string to the buffer, escaping special characters
// as needed for JSON. Valid UTF-8 is passed through as-is. Invalid UTF-8
// sequences are replaced with \ufffd. The string is enclosed in double quotes.
func (buf *Buffer) WriteEscapedString(str string) {
	buf.WriteBytes('"')
	// last is the last index of the string that has been written to the buffer.
	// cur is the current index of the string being processed.
	//
	// Read the string byte by byte and escape any characters that need it.
	// Valid multi-byte UTF-8 sequences are skipped over and included in the
	// next batch flush. Write to the buffer as we go.
	last := 0
	for cur := 0; cur < len(str); {
		b := str[cur]
		if b >= 0x80 {
			_, size := utf8.DecodeRuneInString(str[cur:])
			if size == 1 {
				// \uFFFD is the replacement character for invalid UTF-8
				// sequences (�).
				if last < cur {
					buf.WriteString(str[last:cur])
				}
				buf.WriteString(`\ufffd`)
				cur++
				last = cur
			} else {
				cur += size
			}
		} else if b < 0x20 || b == '"' || b == '\\' {
			if last < cur {
				buf.WriteString(str[last:cur])
			}
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
			}
			cur++
			last = cur
		} else {
			cur++
		}
	}
	// Flush remaining characters that don't need escaping
	if last < len(str) {
		buf.WriteString(str[last:])
	}
	buf.WriteBytes('"')
}

// WriteBase64 writes a byte slice to the buffer as a base64-encoded string.
func (buf *Buffer) WriteBase64(b64enc *base64.Encoding, data []byte) {
	buf.WriteBytes('"')
	buf.b = b64enc.AppendEncode(buf.b, data)
	buf.WriteBytes('"')
}

//
// Buffer pool
//

// Buffers are pooled to reduce allocations.
var bufferPool = sync.Pool{
	New: func() any {
		return &Buffer{make([]byte, 0, bufferSize)}
	},
}

func getBuffer() *Buffer {
	buf, _ := bufferPool.Get().(*Buffer)
	return buf
}

func putBuffer(buf *Buffer) {
	const maxCap = 10 * bufferSize
	if cap(buf.b) > maxCap {
		// If the buffer is too large, let it get garbage collected.
		// This avoids keeping large buffers in the pool to reduce memory usage.
		return
	}
	buf.b = buf.b[:0] // Reset the underlying slice
	bufferPool.Put(buf)
}
