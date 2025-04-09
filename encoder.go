package blip

// Encoder is an interface for encoding log messages.
type Encoder interface {
	// EncodeTime encodes the time of the log message.
	EncodeTime(buf *Buffer, ts string)
	// EncodeLevel encodes the log level of the message.
	EncodeLevel(buf *Buffer, lev Level)
	// EncodeMessage encodes the log message.
	EncodeMessage(buf *Buffer, msg string)
	// EncodeFields encodes the fields of the log message.
	EncodeFields(buf *Buffer, lev Level, fields *[]Field)
	// EncodeStackTrace encodes the stack trace of the log message.
	EncodeStackTrace(buf *Buffer, skip int)
}
