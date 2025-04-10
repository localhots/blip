package blip_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/localhots/blip"
	"github.com/localhots/blip/ctx/log"
)

//
// Fuzz
//

func FuzzBlip(f *testing.F) {
	cfg := blip.DefaultConfig()
	cfg.Encoder = blip.NewJSONEncoder()
	ctx := context.Background()

	// Seed inputs
	f.Add("test message with key-value pair", "key", "value")
	f.Add("test message with spaces in value", "key", "value with spaces")
	// Different alphabets
	f.Add("message with special chars", "yay", "!@#$%^&*()")
	f.Add("消息带有Unicode字符", "键", "你好")
	f.Add("メッセージにUnicode文字が含まれています", "キー", "こんにちは")
	f.Add("Сообщение с Unicode символами", "ключ", "Привет")
	f.Add("رسالة تحتوي على أحرف Unicode", "مفتاح", "مرحبا")
	f.Add("Besked med Unicode-tegn", "nøgle", "Hej")
	// Invalid UTF-8 sequences
	f.Add("message with invalid UTF-8 chars \xff\xfe\xfd", "key", string([]byte{0xff, 0xfe, 0xfd}))
	f.Add("message with invalid UTF-8 chars \x80\x81\x82", "key", string([]byte{0x80, 0x81, 0x82}))
	f.Add("message with invalid UTF-8 chars \xED\xA0\x80", "key", string([]byte{0xED, 0xA0, 0x80}))
	f.Add("message with invalid UTF-8 chars \xED\xB0\x80", "key", string([]byte{0xED, 0xB0, 0x80}))
	f.Add("message with invalid UTF-8 sequence \xED\xA0\x80\xED\xB0\x80", "key", string([]byte{0xED, 0xA0, 0x80, 0xED, 0xB0, 0x80}))
	f.Add("message with repeated invalid UTF-8 sequence \xED\xA0\x80\xED\xB0\x80\xED\xA0\x80", "key", string([]byte{0xED, 0xA0, 0x80, 0xED, 0xB0, 0x80, 0xED, 0xA0, 0x80}))
	f.Add("message with long invalid UTF-8 sequence \xED\xA0\x80\xED\xB0\x80\xED\xA0\x80\xED\xB0\x80", "key", string([]byte{0xED, 0xA0, 0x80, 0xED, 0xB0, 0x80, 0xED, 0xA0, 0x80, 0xED, 0xB0, 0x80}))
	// Very long strings
	f.Add("message with long string "+strings.Repeat("a", 1000), "key", "value with a long string "+strings.Repeat("a", 1000))
	f.Add("message with very long string "+strings.Repeat("a", 10000), "key", "value with a very long string "+strings.Repeat("a", 10000))

	f.Fuzz(func(t *testing.T, msg string, fkey string, fval string) {
		var buf bytes.Buffer
		cfg.Output = &buf
		logger := blip.New(cfg)
		logger.Info(ctx, msg, log.F{fkey: fval})

		// Validate the output
		if buf.Len() == 0 {
			t.Error("Expected non-empty buffer")
		}
		var out any
		if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
			t.Logf("msg=%q key=%q val=%q", msg, fkey, fval)
			t.Log(buf.String())
			t.Errorf("Failed to unmarshal JSON: %v", err)
		}
		if _, ok := out.(map[string]any); !ok {
			t.Errorf("Expected JSON object, got %T", out)
		}
	})
}

//
// Benchmarks
//

func BenchmarkJSON(b *testing.B) {
	log.Setup(blip.Config{
		Level:           blip.LevelDebug,
		Output:          io.Discard,
		StackTraceLevel: blip.LevelError,
		Encoder:         blip.NewJSONEncoder(),
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
		Level:         blip.LevelDebug,
		Output:        io.Discard,
		Time:          true,
		TimeFormat:    "2006-01-02 15:04:05.000",
		TimePrecision: 1 * time.Millisecond,
		Encoder: blip.ConsoleEncoder{
			Color:           false,
			MinMessageWidth: 0,
			SortFields:      false,
		},
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
		Level:         blip.LevelDebug,
		Output:        io.Discard,
		Time:          true,
		TimeFormat:    "2006-01-02 15:04:05.000",
		TimePrecision: 1 * time.Millisecond,
		Encoder: blip.ConsoleEncoder{
			Color:           true,
			MinMessageWidth: 40,
			SortFields:      false,
		},
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
		Level:         blip.LevelDebug,
		Output:        io.Discard,
		Time:          true,
		TimeFormat:    "2006-01-02 15:04:05.000",
		TimePrecision: 1 * time.Millisecond,
		Encoder: blip.ConsoleEncoder{
			Color:           true,
			MinMessageWidth: 40,
			SortFields:      true,
		},
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
		Level:         blip.LevelDebug,
		Output:        io.Discard,
		Time:          true,
		TimeFormat:    "2006-01-02 15:04:05.000",
		TimePrecision: 1 * time.Millisecond,
		Encoder: blip.ConsoleEncoder{
			Color:           true,
			MinMessageWidth: 40,
			SortFields:      true,
		},
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
