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

func Trace(msg string, fields ...F) {
	logger.Trace(context.Background(), msg, fields...)
}

func Debug(msg string, fields ...F) {
	logger.Debug(context.Background(), msg, fields...)
}

func Info(msg string, fields ...F) {
	logger.Info(context.Background(), msg, fields...)
}

func Warn(msg string, fields ...F) {
	logger.Warn(context.Background(), msg, fields...)
}

func Error(msg string, fields ...F) {
	logger.Error(context.Background(), msg, fields...)
}

func Panic(msg string, fields ...F) {
	logger.Panic(context.Background(), msg, fields...)
}

func Fatal(msg string, fields ...F) {
	logger.Fatal(context.Background(), msg, fields...)
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
