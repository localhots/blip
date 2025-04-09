// Package log provides a shim for the blip logger, allowing to use it without
// context. Use it directly or copy into your project to tailor it to your
// needs.
package log

import (
	"context"

	"github.com/localhots/blip"
)

var logger *blip.Logger

// F is a field set used to log key-value pairs.
type F = blip.F

// Setup initializes the logger with the given configuration.
func Setup(cfg blip.Config) {
	logger = blip.New(cfg)
}

// Trace is used to log a message at the Trace level.
func Trace(msg string, fields ...F) {
	logger.Trace(context.Background(), msg, fields...)
}

// Debug is used to log a message at the Debug level.
func Debug(msg string, fields ...F) {
	logger.Debug(context.Background(), msg, fields...)
}

// Info is used to log a message at the Info level.
func Info(msg string, fields ...F) {
	logger.Info(context.Background(), msg, fields...)
}

// Warn is used to log a message at the Warn level.
func Warn(msg string, fields ...F) {
	logger.Warn(context.Background(), msg, fields...)
}

// Error is used to log a message at the Error level.
func Error(msg string, fields ...F) {
	logger.Error(context.Background(), msg, fields...)
}

// Panic is used to log a message at the Panic level.
func Panic(msg string, fields ...F) {
	logger.Panic(context.Background(), msg, fields...)
}

// Fatal is used to log a message at the Fatal level and exit the program.
func Fatal(msg string, fields ...F) {
	logger.Fatal(context.Background(), msg, fields...)
}

// Cause returns a field set that wraps the given error in a standardized way.
func Cause(err error) F {
	return F{"error": err.Error()}
}

// WithContext adds logging fields to the context.
func WithContext(ctx context.Context, fields F) context.Context {
	return blip.WithContext(ctx, fields)
}

// FromContext retrieves the field set from the context.
func FromContext(ctx context.Context) F {
	return blip.FromContext(ctx)
}
