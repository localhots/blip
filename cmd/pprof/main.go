// Package main enables CPU profiling and logs 10M entries.
package main

import (
	"context"
	"flag"
	"io"
	"os"
	"runtime/pprof"

	"github.com/localhots/blip"
	"github.com/localhots/blip/ctx/log"
)

func main() {
	cpuprofile := flag.String("cpuprofile", "", "Write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	log.Setup(blip.Config{
		Level:           blip.LevelDebug,
		Output:          io.Discard,
		StackTraceLevel: blip.LevelError,
		Encoder:         blip.NewJSONEncoder(),
	})
	ctx := context.Background()

	for range 10_000_000 {
		log.Info(ctx, "Starting task", log.F{
			"device_unique_id": "G4000E-1000-F",
			"task_id":          123456,
			"status":           "success",
			"template_name":    "index.tpl",
		})
	}
}
