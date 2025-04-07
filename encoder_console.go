package blip

type ConsoleEncoder struct {
	TimeFormat      string
	MinMessageWidth int
	SortFields      bool
	Color           bool
}

var _ Encoder = (*ConsoleEncoder)(nil)

func NewConsoleEncoder(cfg Config) ConsoleEncoder {
	return ConsoleEncoder{
		TimeFormat:      cfg.TimeFormat,
		MinMessageWidth: cfg.MinMessageWidth,
		SortFields:      cfg.SortFields,
		Color:           cfg.Color,
	}
}

func (e ConsoleEncoder) EncodeTime(buf *Buffer, ts string) {
	buf.WriteString(ts)
	buf.WriteBytes(' ')
}

func (e ConsoleEncoder) EncodeLevel(buf *Buffer, lev Level) {
	e.writeColorized(buf, lev, lev.String())
	buf.WriteBytes(' ')
}

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
		buf.WriteAny(f.Value)
	}
}

func (e ConsoleEncoder) EncodeStackTrace(buf *Buffer, lev Level, skip int) {
	// Print stack trace but skip the first 4 frames which are part of the
	// logger itself.
	buf.WriteString(stackTrace(skip))
	buf.WriteBytes('\n')
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
