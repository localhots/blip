package blip

import (
	"context"
	"maps"
)

type contextKey struct{}

func WithContext(ctx context.Context, fields F) context.Context {
	existing := FromContext(ctx)
	if existing == nil {
		existing = fields
	} else {
		maps.Copy(existing, fields)
	}
	return context.WithValue(ctx, contextKey{}, existing)
}

func FromContext(ctx context.Context) F {
	if v, ok := ctx.Value(contextKey{}).(F); ok {
		return v
	}
	return nil
}
