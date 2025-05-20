package blip

import (
	"context"
	"maps"
)

type contextKey struct{}

// ContextWithFields adds fields to the context. If the context already has
// fields, it merges the new fields with the existing ones.
func ContextWithFields(ctx context.Context, fields F) context.Context {
	existing := FieldsFromContext(ctx)
	if existing == nil {
		existing = fields
	} else {
		maps.Copy(existing, fields)
	}
	return context.WithValue(ctx, contextKey{}, existing)
}

// FieldsFromContext retrieves fields from the context. If no fields are found,
// it returns nil.
func FieldsFromContext(ctx context.Context) F {
	if v, ok := ctx.Value(contextKey{}).(F); ok {
		return v
	}
	return nil
}
