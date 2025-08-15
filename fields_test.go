package blip

import (
	"context"
	"slices"
	"testing"
)

func TestMakeFields(t *testing.T) {
	ctx := context.Background()
	ctx = ContextWithFields(ctx, F{
		"a": 1,
		"b": 2,
	})
	fields := makeFields(ctx, []F{
		{"c": 3, "a": -1},
		{"c": 4},
	})
	sortFields(*fields)
	defer putFields(fields)

	exp := []Field{
		{"a", -1},
		{"b", 2},
		{"c", 4},
	}
	if !slices.Equal(exp, *fields) {
		t.Errorf("expected %v, got %v", exp, *fields)
	}
}

func TestMakeFieldsEmpty(t *testing.T) {
	ctx := context.Background()
	fields := makeFields(ctx, []F{})
	if fields != nil {
		t.Errorf("expected nil, got %v", fields)
	}
}

func TestMakeFieldsOnlyContext(t *testing.T) {
	ctx := context.Background()
	ctx = ContextWithFields(ctx, F{
		"a": 1,
	})
	fields := makeFields(ctx, nil)
	if fields == nil {
		t.Fatal("expected non-nil fields")
	}
	defer putFields(fields)

	exp := []Field{
		{"a", 1},
	}
	if !slices.Equal(exp, *fields) {
		t.Errorf("expected %v, got %v", exp, *fields)
	}
}

func TestSortFields(t *testing.T) {
	fields := []Field{
		{"b", 2},
		{"a", 1},
		{"d", 5},
		{"f", 5},
		{"d", 5},
	}
	sortFields(fields)
	exp := []Field{
		{"a", 1},
		{"b", 2},
		{"d", 5},
		{"d", 5},
		{"f", 5},
	}
	if !slices.Equal(exp, fields) {
		t.Errorf("expected %v, got %v", exp, fields)
	}
}
