package blip_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/localhots/blip"
	"github.com/localhots/blip/ctx/log"
)

func TestLogger(t *testing.T) {
	cfg := blip.DefaultConfig()
	cfg.Level = blip.LevelDebug
	log.Setup(cfg)
	ctx := context.Background()
	err := errors.New("task already exists")

	log.Debug(ctx, "Message without fields")
	log.Debug(ctx, "Callback received", log.F{
		"device_unique_id": "G4000E-1000-F",
		"task_id":          123456,
		"status":           "success",
		"template_name":    "index.tpl",
	})
	log.Info(ctx, "Starting task", log.F{
		"device_unique_id": "G4000E-1000-F",
		"task_id":          123456,
	})
	log.Info(ctx, "Extremely long message, sorry about that; it is meant to prove that buffer can grow"+strings.Repeat(" @@@", 300), log.F{
		"device_unique_id": "G4000E-1000-F",
		"task_id":          123456,
	})
	log.Warn(ctx, "Duplicate task, but also this message exceeds 40 characters", log.F{
		"task_id": 123456,
	})
	log.Warn(ctx, "Duplicate task but exactly 40 characters", log.F{
		"task_id": 123456,
	})
	log.Warn(blip.WithContext(ctx, log.F{"foo": "bar"}), "Duplicate task is exactly 39 characters", log.F{
		"task_id": 123456,
	})
	log.Error(ctx, "Failed to process task", log.Cause(err), log.F{
		"task_id": 123456,
	})
	log.Fatal(ctx, "Failed to start service", log.Cause(err), log.F{
		"service": "api",
	})
}

func BenchmarkBare(b *testing.B) {
	log.Setup(blip.Config{
		Level:           blip.LevelDebug,
		Output:          io.Discard,
		StackTraceLevel: blip.LevelError,
	})
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task", log.F{
			"device_unique_id": "G4000E-1000-F",
			"task_id":          123456,
			"status":           "success",
			"template_name":    "index.tpl",
		})
	}
}

func BenchmarkOptimized(b *testing.B) {
	log.Setup(blip.Config{
		Level:           blip.LevelDebug,
		Output:          io.Discard,
		Time:            true,
		TimeFormat:      "2006-01-02 15:04:05.000",
		TimePrecision:   1 * time.Millisecond,
		Color:           false,
		MinMessageWidth: 0,
		SortFields:      false,
		StackTraceLevel: blip.LevelError,
	})
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task", log.F{
			"device_unique_id": "G4000E-1000-F",
			"task_id":          123456,
			"status":           "success",
			"template_name":    "index.tpl",
		})
	}
}

func BenchmarkPretty(b *testing.B) {
	log.Setup(blip.Config{
		Level:           blip.LevelDebug,
		Output:          io.Discard,
		Time:            true,
		TimeFormat:      "2006-01-02 15:04:05.000",
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      false,
		StackTraceLevel: blip.LevelError,
	})
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task", log.F{
			"device_unique_id": "G4000E-1000-F",
			"task_id":          123456,
			"status":           "success",
			"template_name":    "index.tpl",
		})
	}
}

func BenchmarkPrettySorted(b *testing.B) {
	log.Setup(blip.Config{
		Level:           blip.LevelDebug,
		Output:          io.Discard,
		Time:            true,
		TimeFormat:      "2006-01-02 15:04:05.000",
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      true,
		StackTraceLevel: blip.LevelError,
	})
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task", log.F{
			"device_unique_id": "G4000E-1000-F",
			"task_id":          123456,
			"status":           "success",
			"template_name":    "index.tpl",
		})
	}
}

func BenchmarkPrettySortedContext(b *testing.B) {
	log.Setup(blip.Config{
		Level:           blip.LevelDebug,
		Output:          io.Discard,
		Time:            true,
		TimeFormat:      "2006-01-02 15:04:05.000",
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      true,
		StackTraceLevel: blip.LevelError,
	})
	ctx := context.Background()
	ctx = blip.WithContext(ctx, log.F{"foo": "bar"})
	ctx = blip.WithContext(ctx, log.F{"one": "two"})

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task", log.F{
			"device_unique_id": "G4000E-1000-F",
			"task_id":          123456,
			"status":           "success",
			"template_name":    "index.tpl",
		})
	}
}
