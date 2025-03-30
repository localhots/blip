package blip

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type buffer struct {
	b []byte
}

const bufferSize = 1024

func (buf *buffer) writeSpace() {
	buf.writeByte(' ')
}

func (buf *buffer) writeNewline() {
	buf.writeByte('\n')
}

//
// Encoding
//

func encodeField(buf *buffer, val any) {
	switch v := val.(type) {
	case string:
		buf.writeString(v)
	case []byte:
		buf.writeBytes(v)
	case int:
		buf.writeInt(int64(v))
	case int8:
		buf.writeInt(int64(v))
	case int16:
		buf.writeInt(int64(v))
	case int32:
		buf.writeInt(int64(v))
	case int64:
		buf.writeInt(v)
	case uint:
		buf.writeUint(uint64(v))
	case uint8:
		buf.writeUint(uint64(v))
	case uint16:
		buf.writeUint(uint64(v))
	case uint32:
		buf.writeUint(uint64(v))
	case uint64:
		buf.writeUint(v)
	case float32:
		buf.writeFloat(float64(v), 32)
	case float64:
		buf.writeFloat(v, 64)
	case bool:
		buf.writeBool(v)
	case time.Duration:
		buf.writeDuration(v.Truncate(DurationPrecision))
	case time.Time:
		buf.writeTime(v, TimeFormat)
	default:
		// TODO: Add support for custom encoders
		buf.writeString(fmt.Sprint(v))
	}
}

func (buf *buffer) writeByte(b byte) {
	buf.b = append(buf.b, b)
}

func (buf *buffer) writeBytes(b []byte) {
	buf.b = append(buf.b, b...)
}

func (buf *buffer) writeString(str string) {
	buf.b = append(buf.b, str...)
}

func (buf *buffer) writeInt(i int64) {
	buf.b = strconv.AppendInt(buf.b, i, 10)
}

func (buf *buffer) writeUint(i uint64) {
	buf.b = strconv.AppendUint(buf.b, i, 10)
}

func (buf *buffer) writeFloat(f float64, bitSize int) {
	buf.b = strconv.AppendFloat(buf.b, f, 'f', -1, bitSize)
}

func (buf *buffer) writeBool(b bool) {
	buf.b = strconv.AppendBool(buf.b, b)
}

func (buf *buffer) writeDuration(d time.Duration) {
	buf.b = append(buf.b, d.String()...)
}

func (buf *buffer) writeTime(t time.Time, format string) {
	buf.b = t.AppendFormat(buf.b, format)
}

//
// Buffer pool
//

// Buffers are pooled to reduce allocations.
var bufferPool = sync.Pool{
	New: func() any {
		buf := buffer{make([]byte, 0, bufferSize)}
		return &buf
	},
}

func getBuffer() *buffer {
	s, _ := bufferPool.Get().(*buffer)
	s.b = s.b[:0] // Reset the underlying slice
	return s
}

func putBuffer(buf *buffer) {
	bufferPool.Put(buf)
}
