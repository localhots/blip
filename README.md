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

![Demo](https://github.com/user-attachments/assets/55175d0b-80a5-4fb9-9088-a331a6f3e372)

Blip does not provide Printf-like methods, instead it encourages the use of
fields. Fields are defined as a map, making it look nicely indented with `gofmt`.
There is also a standardized helper for the error type: `log.Cause(err)`.

```go
log.Error("Failed to process task", log.Cause(err), log.F{
	"task_id": 123456,
})
```

The use of `map[string]any` to define fields is optimized by the compiler and
avoids stressing the garbage collector thanks to memory pooling, making it an
ergonomic and worry-free way to log values without concern for their types.

## Context

Fields can be added to context allowing them to be propagated along with it.
Such fields will be logged with every message.

```go
ctx := context.Background()
ctx = log.ContextWithFields(ctx, log.F{
	"task_id": task.ID,
})

if err := runTask(ctx, task); err != nil {
	log.Error(ctx, "Task failed", log.Cause(err))
}
```

## Use

Blip offers both an
[instance-based API](https://pkg.go.dev/github.com/localhots/blip) and a
package-level API. In fact, two package-level variants:
[one with context](https://pkg.go.dev/github.com/localhots/blip/ctx/log),
[one without](https://pkg.go.dev/github.com/localhots/blip/noctx/log).
These can be used directly or copied into a project as a foundation for building
a custom logging package.

### Instance Based

```go
import "github.com/localhots/blip"

ctx := context.Background()
logger := blip.New(blip.DefaultConfig())
logger.Info(ctx, "Callback received", log.F{
	"task_id": 123456,
})
```

### With Context

```go
import "github.com/localhots/blip/ctx/log"

ctx := context.Background()
log.Info(ctx, "Callback received", log.F{
	"task_id": 123456,
})
```

### Without Context

```go
import "github.com/localhots/blip/noctx/log"

log.Info("Callback received", log.F{
	"task_id": 123456,
})
```

## Configuration

The logger can be configured with:

- `Level` — minimum logging level (`Info` by default)
- `Output` — log destination (`stderr` by default)
- `Encoder` — console, JSON, or a custom encoder (console by default)
- `StackTraceLevel` — minimum level at which stack traces are logged (`Error` by default)

Blip includes two built-in encoders: console and JSON, both are further
customizable.

### Console Encoder

- `TimeFormat`
- `TimePrecision` — when positive, caches timestamps until they change by the
  given amount
- `MinMessageWidth` — controls padding between the message and fields
- `SortFields` — enables sorting of fields
- `Color` — enables color and bold text for messages

Fields are sorted using insertion sort, which is highly efficient for small
collections.

### JSON Encoder

- `TimeFormat`
- `TimePrecision` — same behavior as in the console encoder
- `Base64Encoding` — customizes how byte slices are base64-encoded
- `KeyTime`, `KeyLevel`, `KeyMsg`, `KeyStackTrace` — controls the JSON keys for
  corresponding values

## Performance

Blip makes a few intentional compromises in favor of ergonomics and developer
productivity: fields are passed as a map, and reflection is used instead of
specialized functions to log field values. However, it mitigates most drawbacks
through memory pooling for both messages and fields, and by caching timestamps
to avoid expensive time serialization.

Overall, if developer productivity matters more than chasing the absolute lowest
latency, and a slight performance tradeoff per log call isn't that critical,
Blip might be the right choice.

## Comparison to Other Loggers

While it's impossible to make a perfectly fair comparison, here are a few notes
on how Blip fares against other popular logging libraries.

The main motivation behind Blip was to build a logger that is nicer than Logrus
and is as fast as Zerolog.

### [log/slog](https://pkg.go.dev/log/slog)

Slog is a structured logger introduced in Go 1.21. It accepts fields as variadic
`any` arguments or as typed attributes, and supports both console and JSON
encoders. It's allocation-free, fast (still ~2x slower than Blip) but not
particularly pretty.

### [zerolog](https://github.com/rs/zerolog)

Zerolog is an excellent logger and among the fastest available. Its performance
comes from using typed functions to provide fields efficiently. When used as
intended, Zerolog is roughly 25% faster than Blip.

However, if the `Any` method is used instead of `Str`, `Int`, and other typed
functions, it starts allocating memory and falls behind, performing about twice
as slow as Blip.

In pretty mode, Zerolog first encodes messages as JSON, then parses and
re-formats them for console output, which absolutely tanks its performance
compared to the competition. Clearly an afterthought.

### [zap](https://github.com/uber-go/zap)

Zap is reasonably fast but doesn't offer a true pretty mode. It achieves most of
its performance gains through message sampling. With sampling disabled, its
performance drops and it becomes much slower than Blip in all use cases.

### [logrus](https://github.com/sirupsen/logrus)

Logrus offers some of the nicest console formatting and was a major inspiration
for Blip's console encoder. Unfortunately, it doesn't perform well, and active
development has effectively stopped.

### [phuslog](https://github.com/phuslu/log)

Phuslog appears to be the absolute fastest logger available. It employs a lot of
complex low-level code to achieve its performance, demonstrating the significant
effort that went into it.

One reason for its speed is its highly optimized time serialization.
Unfortunately, it accomplishes this by linking to unexported functions from the
`time` package, which has
[caused some grief](https://github.com/golang/go/blob/dad4f39971d89b56224d1eb44121305b1c0ef711/src/time/time.go#L1304-L1314)
for the Go language developers.
