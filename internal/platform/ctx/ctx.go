// Package commonctx provides type-safe utilities for managing custom values
// within a context.Context using generics and unique context keys.
package commonctx

import "context"

// CtxKey represents a unique identifier used as a key in a context.
// It contains a descriptive name to assist with debugging.
type CtxKey struct {
	Name string
}

// SetValue returns a copy of the parent context in which the value associated with key is value.
// The function uses generics to ensure type safety for the provided value.
func SetValue[T any](ctx context.Context, key *CtxKey, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetValue extracts a value of type T associated with key from the context.
// It returns the zero value of T and false if the key is not present or the value type is incorrect.
//
//nolint:ireturn
func GetValue[T any](ctx context.Context, key *CtxKey) (T, bool) {
	value, ok := ctx.Value(key).(T)

	return value, ok
}
