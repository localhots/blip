package blip

import (
	"context"
	"maps"
)

type contextKey struct{}

// WithContext adds fields to the context. If the context already has fields,
// it merges the new fields with the existing ones.
func WithContext(ctx context.Context, fields F) context.Context {
	existing := FromContext(ctx)
	if existing == nil {
		existing = fields
	} else {
		maps.Copy(existing, fields)
	}
	return context.WithValue(ctx, contextKey{}, existing)
}

// FromContext retrieves fields from the context. If no fields are found,
// it returns nil.
func FromContext(ctx context.Context) F {
	if v, ok := ctx.Value(contextKey{}).(F); ok {
		return v
	}
	return nil
}
