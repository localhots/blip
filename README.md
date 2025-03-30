# Blip

Blip is a logger.

```go
import "github.com/localhots/blip"
import "github.com/localhots/blip/ctx/log"

log.Setup(blip.DefaultConfig())

log.Debug(ctx, "Callback received", log.F{
	"device_unique_id": "G4000E-1000-F",
	"task_id":          123456,
	"status":           "success",
	"template_name":    "index.tpl",
})
```
