package memory_test

import (
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/platform/caching/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewLRUCache_CapacityClamped(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[int](0)

	cache.Set("a", 1)
	cache.Set("b", 2)

	assert.Equal(t, 1, cache.Len())
}

func TestLRUCache_SetGet_Basic(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[int](10)

	vMissing, okMissing := cache.Get("missing")
	assert.False(t, okMissing)
	assert.Equal(t, 0, vMissing)
	assert.Equal(t, 0, cache.Len())

	cache.Set("a", 1)
	cache.Set("b", 2)

	assert.Equal(t, 2, cache.Len())

	vA, okA := cache.Get("a")
	assert.True(t, okA)
	assert.Equal(t, 1, vA)

	vB, okB := cache.Get("b")
	assert.True(t, okB)
	assert.Equal(t, 2, vB)
}

func TestLRUCache_Set_UpdatesExisting(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[int](10)

	cache.Set("key", 1)
	assert.Equal(t, 1, cache.Len())

	v1, ok1 := cache.Get("key")
	assert.True(t, ok1)
	assert.Equal(t, 1, v1)

	cache.Set("key", 42)

	v2, ok2 := cache.Get("key")
	assert.True(t, ok2)
	assert.Equal(t, 42, v2)

	assert.Equal(t, 1, cache.Len())
}

func TestLRUCache_Eviction_WithoutAccess(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[string](2)

	cache.Set("a", "A")
	cache.Set("b", "B")
	assert.Equal(t, 2, cache.Len())

	cache.Set("c", "C")
	assert.Equal(t, 2, cache.Len())

	_, okA := cache.Get("a")
	valB, okB := cache.Get("b")
	valC, okC := cache.Get("c")

	assert.False(t, okA, "key 'a' should be evicted as LRU")
	assert.True(t, okB, "key 'b' should remain in cache")
	assert.True(t, okC, "key 'c' should be present in cache")
	assert.Equal(t, "B", valB)
	assert.Equal(t, "C", valC)
}

func TestLRUCache_Eviction_WithAccess(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[string](2)

	cache.Set("a", "A")
	cache.Set("b", "B")
	assert.Equal(t, 2, cache.Len())

	_, _ = cache.Get("a")

	cache.Set("c", "C")
	assert.Equal(t, 2, cache.Len())

	_, okA := cache.Get("a")
	_, okB := cache.Get("b")
	_, okC := cache.Get("c")

	assert.True(t, okA, "key 'a' should remain (recently used)")
	assert.False(t, okB, "key 'b' should be evicted as LRU")
	assert.True(t, okC, "key 'c' should be present in cache")
}

func TestLRUCache_GetAndRemove(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[int](10)

	cache.Set("x", 100)
	assert.Equal(t, 1, cache.Len())

	v, ok := cache.GetAndRemove("x")
	assert.True(t, ok)
	assert.Equal(t, 100, v)
	assert.Equal(t, 0, cache.Len())

	v2, ok2 := cache.Get("x")
	assert.False(t, ok2)
	assert.Equal(t, 0, v2)
	assert.Equal(t, 0, cache.Len())

	v3, ok3 := cache.GetAndRemove("x")
	assert.False(t, ok3)
	assert.Equal(t, 0, v3)
	assert.Equal(t, 0, cache.Len())
}

func TestLRUCache_Iter_Order_MRUToLRU(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[string](10)

	cache.Set("a", "A")
	cache.Set("b", "B")
	cache.Set("c", "C")

	_, _ = cache.Get("a")

	var order []string
	cache.Iter(func(key, _ string) {
		order = append(order, key)
	})

	assert.Equal(t, []string{"a", "c", "b"}, order)
}

func TestLRUCache_IterFromLast_Order_LRUToMRU(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[string](10)

	cache.Set("a", "A")
	cache.Set("b", "B")
	cache.Set("c", "C")

	_, _ = cache.Get("a")

	var order []string
	cache.IterFromLast(func(key, _ string) {
		order = append(order, key)
	})

	assert.Equal(t, []string{"b", "c", "a"}, order)
}

func TestLRUCache_Iter_NilAction_NoPanic(t *testing.T) {
	t.Parallel()

	cache := memory.NewLRUCache[int](10)
	cache.Set("a", 1)
	cache.Set("b", 2)

	cache.Iter(nil)
	cache.IterFromLast(nil)
}
