package blip

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	cfg       Config
	timeCache func(time.Time) string
	lock      sync.Mutex
}

type Config struct {
	Level           Level
	Output          io.Writer
	Time            bool
	TimeFormat      string
	TimePrecision   time.Duration
	Color           bool
	MinMessageWidth int
	SortFields      bool
	StackTraceLevel Level
	StackTraceSkip  int
}

type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

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

var (
	defaultMessageWidth = 40 // characters
	defaultTimeFormat   = "2006-01-02 15:04:05.000"

	DurationPrecision = time.Millisecond
	TimeFormat        = time.RFC3339
)

func New(cfg Config) *Logger {
	l := &Logger{cfg: cfg}
	if l.cfg.TimePrecision > 0 {
		l.timeCache = timeCache(l.cfg.TimeFormat, l.cfg.TimePrecision)
	}
	return l
}

func DefaultConfig() Config {
	return Config{
		Level:           LevelInfo,
		Output:          os.Stderr,
		Time:            true,
		TimeFormat:      defaultTimeFormat,
		TimePrecision:   0, // Disable time cache
		Color:           true,
		MinMessageWidth: defaultMessageWidth,
		SortFields:      true,
		StackTraceLevel: LevelError,
		StackTraceSkip:  4,
	}
}

func (l *Logger) Trace(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level == LevelTrace {
		l.print(LevelTrace, msg, makeFields(ctx, fields))
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelDebug {
		l.print(LevelDebug, msg, makeFields(ctx, fields))
	}
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelInfo {
		l.print(LevelInfo, msg, makeFields(ctx, fields))
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelWarn {
		l.print(LevelWarn, msg, makeFields(ctx, fields))
	}
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelError {
		l.print(LevelError, msg, makeFields(ctx, fields))
	}
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelPanic {
		l.print(LevelPanic, msg, makeFields(ctx, fields))
	}
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...F) {
	l.print(LevelFatal, msg, makeFields(ctx, fields))
	os.Exit(1)
}

//
// Printing
//

func (l *Logger) print(lev Level, msg string, fields *[]Field) {
	buf := getBuffer()
	defer putBuffer(buf)

	l.printTime(buf)
	l.printLevel(buf, lev)
	l.printMessage(buf, msg, len(*fields) > 0)
	l.printFields(buf, lev, *fields)
	putFields(fields)
	buf.writeNewline()
	l.printStackTrace(buf, lev)

	l.lock.Lock()
	l.cfg.Output.Write(buf.b)
	l.lock.Unlock()
}

func (l *Logger) printTime(buf *buffer) {
	if !l.cfg.Time {
		return
	}

	t := time.Now()
	if l.timeCache != nil {
		buf.writeString(l.timeCache(t))
	} else {
		buf.writeTime(t, l.cfg.TimeFormat)
	}
	buf.writeSpace()
}

func (l *Logger) printLevel(buf *buffer, lev Level) {
	l.writeColorized(buf, lev, l.levelName(lev))
	buf.writeSpace()
}

func (l *Logger) printMessage(buf *buffer, msg string, needsPad bool) {
	buf.writeString(msg)
	if l.cfg.MinMessageWidth > 0 {
		// Pad the message to the configured width +2 spaces to separate it from
		// the fields.
		for range l.cfg.MinMessageWidth + 2 - len(msg) {
			buf.writeSpace()
		}
		// If the message is long enough not to be padded, add an extra space to
		// separate it from the fields
		if len(msg) > l.cfg.MinMessageWidth {
			// Separate message from fields with 2 spaces
			buf.writeSpace()
			buf.writeSpace()
		}
	} else if needsPad {
		// Separate message from fields with 2 spaces
		buf.writeSpace()
		buf.writeSpace()
	}
}

func (l *Logger) printFields(buf *buffer, lev Level, fields []Field) {
	if l.cfg.SortFields {
		sortFields(fields)
	}

	for i, f := range fields {
		if i > 0 {
			buf.writeSpace()
		}
		l.writeColorized(buf, lev, f.Key)
		buf.writeByte('=')
		encodeField(buf, f.Value)
	}
}

func (l *Logger) printStackTrace(buf *buffer, lev Level) {
	if lev >= l.cfg.StackTraceLevel {
		// Print stack trace but skip the first 4 frames which are part of the
		// logger itself.
		buf.writeString(stackTrace(l.cfg.StackTraceSkip))
		buf.writeNewline()
	}
}

//
// Helpers
//

func (l *Logger) writeColorized(buf *buffer, lev Level, str string) {
	if !l.cfg.Color {
		buf.writeString(str)
		return
	}

	switch lev {
	case LevelTrace, LevelDebug:
		buf.writeString(colorOffWhite)
	case LevelInfo:
		buf.writeString(colorCyan)
	case LevelWarn:
		buf.writeString(colorYellow)
	case LevelError:
		buf.writeString(colorRed)
	case LevelPanic, LevelFatal:
		buf.writeString(colorRedBg)
		buf.writeString(colorWhite)
	}
	buf.writeString(str)
	buf.writeString(colorReset)
}

// levelName returns level label that is consistently 4 characters long.
func (l *Logger) levelName(lev Level) string {
	switch lev {
	case LevelTrace:
		return "TRAC"
	case LevelDebug:
		return "DEBU"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERRO"
	case LevelPanic:
		return "PANI"
	case LevelFatal:
		return "FATA"
	default:
		panic("unreachable")
	}
}

func stackTrace(skip int) string {
	// Get up to 100 stack frames
	pc := make([]uintptr, 100)
	// +2 frames to skip for runtime.Callers and stackTrace itself
	n := runtime.Callers(skip+2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var buf bytes.Buffer
	for {
		f, more := frames.Next()
		buf.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", f.Function, f.File, f.Line))
		if !more {
			break
		}
	}
	return buf.String()
}
