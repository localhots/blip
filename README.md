# Blip

Blip is a high-performance, generic logging library for Go. It is designed to be
fast, allocation-free, and simple to use, without compromising on performance or
relying on hacks.

```go
log.Info("Callback received", log.F{
	"device_unique_id": "G4000E-1000-F",
	"task_id":          123456,
	"status":           "success",
	"template_name":    "index.tpl",
})
```

![blip_demo](https://github.com/user-attachments/assets/55175d0b-80a5-4fb9-9088-a331a6f3e372)

Blip does not provide Printf-like methods, instead it encourages the use of
fields. Fields are defined as a map, making it look nicely indented with `gofmt`.
There is also a standardized helper for error type: `log.Cause(err)`.

```go
log.Error("Failed to process task", log.Cause(err), log.F{
	"task_id": 123456,
})
```

The use of `map[string]any` to define fields is optimized by the compiler and
avoids stressing the garbage collector thanks to memory pooling, making it an
efficient and worry-free way to log values without concern for their types.

## Context

Fields can be added to context allowing them to be propagated along with it.

```go
ctx := context.Background()
ctx = log.WithContext(ctx, log.F{
	"task_id": task.ID,
})

err := runTask(ctx, task)
```

## Use

Blip offers both an instance-based API and a package-level API. In fact, two
package-level variants: one with context, one without. These can be used
directly or copied into a project as a foundation for building a custom logging
package.

See `ctx/log` and `noctx/log` subpackages for details.

## Comparison to Other Loggers

While it's impossible to make a perfectly fair comparison, here are a few notes
on how Blip compares to other popular logging libraries.

The main motivation behind Blip was to build a logger that's nicer than Logrus
and faster than Zerolog.

### [log/slog](https://pkg.go.dev/log/slog)

Slog is a structured logger introduced in Go 1.21. It accepts fields as variadic
`any` arguments and offers both console and JSON encoders. However, it is
neither particularly fast nor pretty.

### [zerolog](https://github.com/rs/zerolog)

Zerolog is an excellent logger and among the fastest available. Its performance
comes from using typed functions to provide fields efficiently. When used as
intended, Zerolog is roughly 25% faster than Blip.

However, if the `Any` method is used instead of `Str`, `Int`, and other typed
functions, it starts allocating memory and falls behind, performing about twice
as slow as Blip.

In pretty mode, Zerolog first encodes messages as JSON, then parses and
re-formats them for console output, which absolutely tanks its performance
compared to the competition.

### [zap](https://github.com/uber-go/zap)

Zap achieves most of its performance gains through message sampling. With
sampling disabled, its performance drops and it becomes slower than Blip in all
use cases.

### [Logrus](https://github.com/sirupsen/logrus)

Logrus offers some of the nicest console formatting and was a major inspiration
for Blip's console encoder. Unfortunately, it doesn't perform well, and active
development has effectively stopped.
