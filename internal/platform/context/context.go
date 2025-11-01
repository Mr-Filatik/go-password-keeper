// Package context contains functionality for working with context.
package context

import "context"

// CtxKey is a structure describing keys for working with context variables.
type CtxKey struct {
	// Name - name of the key.
	Name string
}

// WithValue sets the value (string) by key in the context.
//
// Parameters:
//   - ctx context.Context: parent context;
//   - key *CtxKey: key;
//   - value string: value as a string.
func WithValue(ctx context.Context, key *CtxKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetValue gets the value (string) from the context by key.
//
// Parameters:
// - ctx context.Context: context;
// - key *CtxKey: key.
func GetValue(ctx context.Context, key *CtxKey) string {
	value := ctx.Value(key)
	if value == nil {
		return ""
	}

	strValue, ok := value.(string)
	if !ok {
		return ""
	}

	return strValue
}
