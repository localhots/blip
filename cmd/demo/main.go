// Package main is a demo application for the blip logger.
package main

import (
	"context"
	"errors"
	"flag"
	"strings"

	"github.com/localhots/blip"
	"github.com/localhots/blip/ctx/log"
)

func main() {
	cfg := blip.DefaultConfig()
	cfg.Level = blip.LevelDebug
	flag.BoolVar(&cfg.Time, "time", cfg.Time, "Log timestamps")
	color := flag.Bool("color", true, "Colorized output")
	sort := flag.Bool("sort", true, "Sort fields")
	width := flag.Int("width", 40, "Min message width")
	encoder := flag.String("enc", "console", "Log encoder (json, console)")
	flag.Parse()
	switch *encoder {
	case "json":
		cfg.Encoder = blip.NewJSONEncoder()
	case "console":
		cfg.Encoder = blip.ConsoleEncoder{
			Color:           *color,
			SortFields:      *sort,
			MinMessageWidth: *width,
		}
	default:
		panic("invalid encoder")
	}

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
	log.Panic(ctx, "Failed to start service", log.Cause(err), log.F{
		"service": "api",
	})
}
