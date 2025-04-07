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
	enc       Encoder
	timeCache func(time.Time) string
	lock      sync.Mutex
}

type Config struct {
	Level           Level
	Output          io.Writer
	Encoder         Encoder
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
	levelInvalid Level = iota
	LevelTrace
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

	timeNow = time.Now
)

func New(cfg Config) *Logger {
	// Set fallback values
	if cfg.Level < LevelTrace || cfg.Level > LevelFatal {
		cfg.Level = LevelInfo
	}
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}
	if cfg.TimeFormat == "" {
		// Assume empty time format is the default one, not disabled
		cfg.TimeFormat = defaultTimeFormat
	}
	if cfg.MinMessageWidth < 0 {
		// Consider negative message width as no padding
		cfg.MinMessageWidth = 0
	}
	if cfg.StackTraceLevel < LevelTrace || cfg.StackTraceLevel > LevelFatal {
		cfg.StackTraceLevel = LevelError
	}
	if cfg.Encoder == nil {
		cfg.Encoder = NewConsoleEncoder(cfg)
	}

	l := &Logger{cfg: cfg, enc: cfg.Encoder}
	if cfg.TimePrecision > 0 {
		l.timeCache = timeCache(l.cfg.TimeFormat, l.cfg.TimePrecision)
	}
	return l
}

func DefaultConfig() Config {
	cfg := Config{
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
	cfg.Encoder = NewConsoleEncoder(cfg)
	return cfg
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

	if l.cfg.Time {
		l.enc.EncodeTime(buf, l.timeString())
	}
	l.enc.EncodeLevel(buf, lev)
	l.enc.EncodeMessage(buf, msg)
	l.enc.EncodeFields(buf, lev, fields)
	if fields != nil {
		putFields(fields)
	}
	if lev >= l.cfg.StackTraceLevel {
		l.enc.EncodeStackTrace(buf, lev, l.cfg.StackTraceSkip)
	}

	l.lock.Lock()
	l.cfg.Output.Write(buf.b)
	l.lock.Unlock()
}

//
// Helpers
//

func (l *Logger) timeString() string {
	if l.timeCache != nil {
		return l.timeCache(timeNow())
	}
	return timeNow().Format(l.cfg.TimeFormat)
}

// levelName returns level label that is consistently 4 characters long.
func (lev Level) String() string {
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
