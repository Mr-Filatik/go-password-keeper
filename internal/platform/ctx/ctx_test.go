package commonctx_test

import (
	"context"
	"testing"

	commonctx "github.com/mr-filatik/go-password-keeper/internal/platform/ctx"
)

func TestContextOperations(t *testing.T) {
	keyString := &commonctx.CtxKey{Name: "test-string"}
	keyWrong := &commonctx.CtxKey{Name: "wrong-key"}

	ctxEmpty := context.Background()
	ctxWithString := commonctx.SetValue(ctxEmpty, keyString, "hello_world")
	ctxWithInt := commonctx.SetValue(ctxEmpty, keyString, 42)

	tests := []struct {
		name          string
		inputCtx      context.Context
		searchKey     *commonctx.CtxKey
		expectedValue string
		expectedOk    bool
	}{
		{
			name:          "successfully get string value",
			inputCtx:      ctxWithString,
			searchKey:     keyString,
			expectedValue: "hello_world",
			expectedOk:    true,
		},
		{
			name:          "return false when key is missing",
			inputCtx:      ctxEmpty,
			searchKey:     keyString,
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "return false when searching with a different key",
			inputCtx:      ctxWithString,
			searchKey:     keyWrong,
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "return false when value type is incorrect (int instead of string)",
			inputCtx:      ctxWithInt,
			searchKey:     keyString,
			expectedValue: "",
			expectedOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := commonctx.GetValue[string](tt.inputCtx, tt.searchKey)

			if ok != tt.expectedOk {
				t.Errorf("GetValue() ok = %v, wantOk %v", ok, tt.expectedOk)
			}

			if result != tt.expectedValue {
				t.Errorf("GetValue() result = %q, want %q", result, tt.expectedValue)
			}
		})
	}
}

func TestContextOperations_PointerUniqueness(t *testing.T) {
	t.Parallel()

	key1 := &commonctx.CtxKey{Name: "user_id"}
	key2 := &commonctx.CtxKey{Name: "user_id"}

	ctx := context.Background()
	ctx = commonctx.SetValue(ctx, key1, "user_one")
	ctx = commonctx.SetValue(ctx, key2, "user_two")

	val1, ok1 := commonctx.GetValue[string](ctx, key1)
	if !ok1 || val1 != "user_one" {
		t.Errorf("failed isolation for key1: got %q, ok: %t", val1, ok1)
	}

	val2, ok2 := commonctx.GetValue[string](ctx, key2)
	if !ok2 || val2 != "user_two" {
		t.Errorf("failed isolation for key2: got %q, ok: %t", val2, ok2)
	}
}
