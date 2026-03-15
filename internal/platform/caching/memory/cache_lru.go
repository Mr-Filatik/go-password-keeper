// Package memory provides functionality for an in-memory cache for an application.
package memory

import (
	"sync"

	"github.com/mr-filatik/go-password-keeper/internal/platform/structure"
)

type lruEntry[T any] struct {
	key   string
	value T
}

// LRUCache is a fixed-capacity least recently used (LRU) cache.
// It stores entries in a map for O(1) lookups and maintains access
// order in a doubly linked list, where the front holds the most
// recently used items and the back holds the least recently used.
// All operations on LRUCache are guarded by mu and are safe for
// concurrent use when accessed through its methods.
type LRUCache[T any] struct {
	mu   sync.RWMutex
	cap  int
	data map[string]*structure.LinkedListNode[lruEntry[T]]
	//nolint:godox
	// TODO: Eliminate double synchronization by replacing the linked list with a non-thread-safe one.
	list *structure.LinkedList[lruEntry[T]]
}

// NewLRUCache returns a new LRUCache with the given capacity.
// If capacity is less than or equal to 1, it is clamped to 1.
func NewLRUCache[T any](capacity int) *LRUCache[T] {
	if capacity <= 1 {
		capacity = 1
	}

	return &LRUCache[T]{
		data: make(map[string]*structure.LinkedListNode[lruEntry[T]]),
		list: structure.NewLinkedList[lruEntry[T]](),
		mu:   sync.RWMutex{},
		cap:  capacity,
	}
}

// Len returns the number of items currently stored in the cache.
// It is safe to call Len concurrently with other operations.
func (c *LRUCache[T]) Len() int {
	c.mu.RLock()

	n := len(c.data)

	c.mu.RUnlock()

	return n
}

// Set stores the given key-value pair in the cache.
// If the key already exists, its value is updated and the entry is
// moved to the front of the list as the most recently used item.
// If the key is new and the cache exceeds its capacity, the least
// recently used item (at the back of the list) is evicted.
// Set is safe for concurrent use.
func (c *LRUCache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.data[key]
	if ok {
		item.Item.value = value
		c.list.MoveFront(item)

		return
	}

	newItem := c.list.PushFront(lruEntry[T]{
		key:   key,
		value: value,
	})
	c.data[key] = newItem

	if len(c.data) > c.cap {
		removeItem := c.list.Last()
		if removeItem != nil {
			c.list.Remove(removeItem)
			delete(c.data, removeItem.Item.key)
		}
	}
}

// Get returns the value associated with the given key and a boolean
// indicating whether the key was found. If the key exists, the
// corresponding entry is moved to the front of the list as the most
// recently used item. Get is safe for concurrent use.
//
//nolint:ireturn
func (c *LRUCache[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T

	item, ok := c.data[key]
	if !ok {
		return zero, false
	}

	c.list.MoveFront(item)

	return item.Item.value, true
}

// GetAndRemove returns the value associated with the given key and a
// boolean indicating whether the key was found, and removes the entry
// from the cache if it exists. The LRU ordering is updated accordingly.
// GetAndRemove is safe for concurrent use.
//
//nolint:ireturn
func (c *LRUCache[T]) GetAndRemove(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T

	item, ok := c.data[key]
	if !ok {
		return zero, false
	}

	c.list.Remove(item)
	delete(c.data, item.Item.key)

	return item.Item.value, true
}

// Iter calls action for each key-value pair in the cache from most
// recently used to least recently used (front to back of the list).
// If action is nil, Iter does nothing.
// The set and order of entries are captured under a read lock before
// action is invoked, so modifications to the cache inside action do
// not affect which entries are visited.
func (c *LRUCache[T]) Iter(action func(key string, value T)) {
	if action == nil {
		return
	}

	c.mu.RLock()

	var entries []lruEntry[T]
	c.list.Iter(func(node *structure.LinkedListNode[lruEntry[T]]) {
		entries = append(entries, node.Item)
	})

	c.mu.RUnlock()

	for _, e := range entries {
		action(e.key, e.value)
	}
}

// IterFromLast calls action for each key-value pair in the cache from
// least recently used to most recently used (back to front of the list).
// If action is nil, IterFromLast does nothing.
// The set and order of entries are captured under a read lock before
// action is invoked, so modifications to the cache inside action do
// not affect which entries are visited.
func (c *LRUCache[T]) IterFromLast(action func(key string, value T)) {
	if action == nil {
		return
	}

	c.mu.RLock()

	var entries []lruEntry[T]
	c.list.IterFromLast(func(node *structure.LinkedListNode[lruEntry[T]]) {
		entries = append(entries, node.Item)
	})

	c.mu.RUnlock()

	for _, e := range entries {
		action(e.key, e.value)
	}
}
