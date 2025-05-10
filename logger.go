// Package blip provides a simple and efficient logging library. It is designed
// to be easy to use and flexible, allowing to customize the logging output
// format, level, and other settings.
// It supports JSON and console output formats and provides an inteface for
// custom encoders.
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

// Logger is a the main structure used to log messages.
type Logger struct {
	cfg  Config
	enc  Encoder
	lock sync.Mutex
}

// Config is the configuration structure for the logger.
type Config struct {
	Level           Level
	Output          io.Writer
	Encoder         Encoder
	StackTraceLevel Level
	StackTraceSkip  int
}

// Level is the log level type.
type Level int

// LevelStringer is a function type that converts a log level to a string.
type LevelStringer func(Level) string

// Supported log levels.
const (
	LevelTrace Level = iota + 1
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

var (
	defaultMessageWidth  = 40 // characters
	defaultTimeFormat    = "2006-01-02 15:04:05.000"
	defaultTimePrecision = 1 * time.Millisecond

	// DurationFieldPrecision controls how duration values are truncated when
	// logged.
	DurationFieldPrecision = time.Millisecond
	// TimeFieldFormat controls the format used for time field values. Log entry
	// timestamps are configured with the TimeFieldFormat field in the Config
	// struct.
	TimeFieldFormat = time.RFC3339

	timeNow = time.Now
)

// New creates a new Logger instance with the given configuration.
func New(cfg Config) *Logger {
	// Set fallback values
	if cfg.Level < LevelTrace || cfg.Level > LevelFatal {
		cfg.Level = LevelInfo
	}
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}
	if cfg.StackTraceLevel < LevelTrace || cfg.StackTraceLevel > LevelFatal {
		cfg.StackTraceLevel = LevelError
	}
	if cfg.Encoder == nil {
		cfg.Encoder = NewConsoleEncoder()
	}

	return &Logger{
		cfg: cfg,
		enc: cfg.Encoder,
	}
}

// DefaultConfig returns a default configuration for the logger.
func DefaultConfig() Config {
	return Config{
		Level:           LevelInfo,
		Output:          os.Stderr,
		StackTraceLevel: LevelError,
		StackTraceSkip:  4,
		Encoder:         NewConsoleEncoder(),
	}
}

// Trace is used to log a message at the Trace level.
func (l *Logger) Trace(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level == LevelTrace {
		l.print(LevelTrace, msg, makeFields(ctx, fields))
	}
}

// Debug is used to log a message at the Debug level.
func (l *Logger) Debug(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelDebug {
		l.print(LevelDebug, msg, makeFields(ctx, fields))
	}
}

// Info is used to log a message at the Info level.
func (l *Logger) Info(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelInfo {
		l.print(LevelInfo, msg, makeFields(ctx, fields))
	}
}

// Warn is used to log a message at the Warn level.
func (l *Logger) Warn(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelWarn {
		l.print(LevelWarn, msg, makeFields(ctx, fields))
	}
}

// Error is used to log a message at the Error level.
func (l *Logger) Error(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelError {
		l.print(LevelError, msg, makeFields(ctx, fields))
	}
}

// Panic is used to log a message at the Panic level.
func (l *Logger) Panic(ctx context.Context, msg string, fields ...F) {
	if l.cfg.Level <= LevelPanic {
		l.print(LevelPanic, msg, makeFields(ctx, fields))
	}
}

// Fatal is used to log a message at the Fatal level and exit the program.
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

	l.enc.Start(buf)
	l.enc.EncodeTime(buf)
	l.enc.EncodeLevel(buf, lev)
	l.enc.EncodeMessage(buf, msg)
	l.enc.EncodeFields(buf, lev, fields)
	if fields != nil {
		putFields(fields)
	}
	if lev >= l.cfg.StackTraceLevel {
		l.enc.EncodeStackTrace(buf, l.cfg.StackTraceSkip)
	}
	l.enc.End(buf)

	l.lock.Lock()
	_, _ = l.cfg.Output.Write(buf.b)
	l.lock.Unlock()
}

//
// Helpers
//

// LevelShortUppercase returns level label that is first 4 characters in
// uppercase.
func LevelShortUppercase(lev Level) string {
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

// LevelFullLowercase returns level label that is all lowercase.
func LevelFullLowercase(lev Level) string {
	switch lev {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelPanic:
		return "panic"
	case LevelFatal:
		return "fatal"
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

func timeCache(format string, precision time.Duration) func(time.Time) string {
	var lastTime time.Time
	var lastTimeStr string

	return func(t time.Time) string {
		if !lastTime.IsZero() && t.Sub(lastTime) < precision {
			return lastTimeStr
		}

		lastTime = t
		lastTimeStr = t.Format(format)
		return lastTimeStr
	}
}
