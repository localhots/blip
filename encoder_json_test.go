package blip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSONEncoder(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultConfig()
	cfg.Output = &buf
	cfg.Level = LevelDebug
	cfg.Encoder = NewJSONEncoder(cfg)

	logger := New(cfg)
	ctx := context.Background()

	logger.Info(ctx, "Starting task", F{
		"device_unique_id": "G4000E-1000-F",
		"task_id":          123456,
	})
	validateAndPrintJSON(t, buf)
}

func TestJSONEncoderNoFields(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultConfig()
	cfg.Output = &buf
	cfg.Level = LevelDebug
	cfg.Encoder = NewJSONEncoder(cfg)

	logger := New(cfg)
	ctx := context.Background()

	logger.Info(ctx, "Starting task")
	validateAndPrintJSON(t, buf)
}

func TestJSONEncoderNoFieldsNoTime(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultConfig()
	cfg.Output = &buf
	cfg.Level = LevelDebug
	cfg.Time = false
	cfg.Encoder = NewJSONEncoder(cfg)

	logger := New(cfg)
	ctx := context.Background()

	logger.Info(ctx, "Starting task")
	validateAndPrintJSON(t, buf)
}

func TestJSONEncoderMulti(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Level = LevelDebug
	cfg.Encoder = NewJSONEncoder(cfg)

	logger := New(cfg)
	ctx := context.Background()

	logger.Info(ctx, "Starting task")
	logger.Info(ctx, "Starting task", F{
		"device_unique_id": "G4000E-1000-F",
		"task_id":          123456,
	})
	logger.Info(ctx, "Starting task", F{
		"device_unique_id": "G4000E-1000-F",
		"task_id":          123456,
	})
}

func validateAndPrintJSON(t *testing.T, buf bytes.Buffer) {
	t.Helper()
	if buf.Len() == 0 {
		t.Error("Expected non-empty buffer")
	}

	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}
	fmt.Print(buf.String())
}
