// Package log provides a shim for the blip logger, allowing to use it with
// context. Use it directly or copy into your project to tailor it to your
// needs.
package log

import (
	"context"

	"github.com/localhots/blip"
)

var logger = blip.New(blip.DefaultConfig())

// F is a field set used to log key-value pairs.
type F = blip.F

// Setup initializes the logger with the given configuration.
func Setup(cfg blip.Config) {
	logger = blip.New(cfg)
}

// Trace is used to log a message at the Trace level.
func Trace(ctx context.Context, msg string, fields ...F) {
	logger.Trace(ctx, msg, fields...)
}

// Debug is used to log a message at the Debug level.
func Debug(ctx context.Context, msg string, fields ...F) {
	logger.Debug(ctx, msg, fields...)
}

// Info is used to log a message at the Info level.
func Info(ctx context.Context, msg string, fields ...F) {
	logger.Info(ctx, msg, fields...)
}

// Warn is used to log a message at the Warn level.
func Warn(ctx context.Context, msg string, fields ...F) {
	logger.Warn(ctx, msg, fields...)
}

// Error is used to log a message at the Error level.
func Error(ctx context.Context, msg string, fields ...F) {
	logger.Error(ctx, msg, fields...)
}

// Panic is used to log a message at the Panic level.
func Panic(ctx context.Context, msg string, fields ...F) {
	logger.Panic(ctx, msg, fields...)
}

// Fatal is used to log a message at the Fatal level and exit the program.
func Fatal(ctx context.Context, msg string, fields ...F) {
	logger.Fatal(ctx, msg, fields...)
}

// Cause returns a field set that wraps the given error in a standardized way.
func Cause(err error) F {
	return F{"error": err.Error()}
}

// WithContext adds logging fields to the context.
func WithContext(ctx context.Context, fields F) context.Context {
	return blip.WithContext(ctx, fields)
}

// FromContext retrieves logging fields from the context.
func FromContext(ctx context.Context) F {
	return blip.FromContext(ctx)
}
