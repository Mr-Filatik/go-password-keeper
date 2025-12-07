// Package caching provides common entities for caching.
package caching

import (
	"context"
	"time"
)

// ICache describes a common interface for key-value storage.
type ICache[T any] interface {
	// SetValue sets the value by key or an error.
	//
	// Parameters:
	//   - ctx context.Context: context;
	//   - key string: key;
	//   - value T: value;
	//   - expiration time.Duration: expiration date.
	SetValue(
		ctx context.Context,
		key string,
		value T,
		expiration time.Duration,
	) error

	// GetValue returns the value by key or an error.
	//
	// Parameters:
	//   - ctx context.Context: context;
	//   - key string: key.
	GetValue(
		ctx context.Context,
		key string,
	) (T, error)
}
