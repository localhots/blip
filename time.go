package blip

import (
	"time"
)

func timeCache(format string, precision time.Duration) func(time.Time) string {
	var lastTime time.Time
	var lastTimeStr string

	return func(t time.Time) string {
		if !lastTime.IsZero() && t.Sub(lastTime) < precision {
			return lastTimeStr
		}

		lastTime = t
		lastTimeStr = t.Format(format)
		return lastTimeStr
	}
}
