package blip

type Encoder interface {
	EncodeTime(buf *Buffer, ts string)
	EncodeLevel(buf *Buffer, lev Level)
	EncodeMessage(buf *Buffer, msg string)
	EncodeFields(buf *Buffer, lev Level, fields *[]Field)
	EncodeStackTrace(buf *Buffer, lev Level, skip int)
}
