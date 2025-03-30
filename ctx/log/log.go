package log

import (
	"context"

	"github.com/localhots/blip"
)

var logger *blip.Logger

type F = blip.F

func Setup(cfg blip.Config) {
	logger = blip.New(cfg)
}

func Trace(ctx context.Context, msg string, fields ...F) {
	logger.Trace(ctx, msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...F) {
	logger.Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...F) {
	logger.Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...F) {
	logger.Warn(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...F) {
	logger.Error(ctx, msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...F) {
	logger.Panic(ctx, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...F) {
	logger.Fatal(ctx, msg, fields...)
}

// Cause returns a field set that wraps the given error in a standardized way.
func Cause(err error) F {
	return F{"error": err.Error()}
}

func WithContext(ctx context.Context, fields F) context.Context {
	return blip.WithContext(ctx, fields)
}

func FromContext(ctx context.Context) F {
	return blip.FromContext(ctx)
}
