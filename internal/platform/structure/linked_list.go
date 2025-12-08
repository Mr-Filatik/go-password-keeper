// Package structure provides generic data structures.
package structure

import (
	"sync"
	"time"
)

// LinkedListNode is a node in a doubly linked list.
// Item holds the stored value, Next and Prev point to the neighboring nodes,
// and ChangedAt records the time of the last update of the node
// (such as insertion or moving it within the list).
type LinkedListNode[T any] struct {
	Item     T
	next     *LinkedListNode[T]
	prev     *LinkedListNode[T]
	UpdateAt time.Time
}

// LinkedList is a generic, doubly linked list.
// All operations on LinkedList are guarded by mu and are safe for
// concurrent use when accessed through its methods.
type LinkedList[T any] struct {
	first *LinkedListNode[T]
	last  *LinkedListNode[T]
	len   int
	mu    sync.RWMutex
}

// NewLinkedList returns a new, empty LinkedList.
// The returned list is safe for concurrent use.
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		first: nil,
		last:  nil,
		len:   0,
		mu:    sync.RWMutex{},
	}
}

// Len returns the number of nodes in the list.
// It is safe to call Len concurrently with other operations on the list.
func (l *LinkedList[T]) Len() int {
	l.mu.RLock()

	n := l.len

	l.mu.RUnlock()

	return n
}

// PushFront inserts a new node with the given value at the front of the list
// and returns the created node. The node's ExpireAt is set to the current
// time in UTC. PushFront is safe for concurrent use.
func (l *LinkedList[T]) PushFront(value T) *LinkedListNode[T] {
	now := time.Now().UTC()

	l.mu.Lock()
	defer l.mu.Unlock()

	item := &LinkedListNode[T]{
		next:     l.first,
		prev:     nil,
		Item:     value,
		UpdateAt: now,
	}

	if l.first == nil {
		l.first = item
		l.last = item
		l.len = 1

		return item
	}

	l.first.prev = item
	l.first = item
	l.len++

	return item
}

// PushBack inserts a new node with the given value at the back of the list
// and returns the created node. The node's ExpireAt is set to the current
// time in UTC. PushBack is safe for concurrent use.
func (l *LinkedList[T]) PushBack(value T) *LinkedListNode[T] {
	now := time.Now().UTC()

	l.mu.Lock()
	defer l.mu.Unlock()

	item := &LinkedListNode[T]{
		next:     nil,
		prev:     l.last,
		Item:     value,
		UpdateAt: now,
	}

	if l.last == nil {
		l.first = item
		l.last = item
		l.len = 1

		return item
	}

	l.last.next = item
	l.last = item
	l.len++

	return item
}

// MoveFront moves the given node to the front of the list.
// If item is nil, the list is empty, or item is already the first node,
// MoveFront does nothing. The node's ExpireAt field is updated to the
// current time in UTC. MoveFront is safe for concurrent use.
func (l *LinkedList[T]) MoveFront(item *LinkedListNode[T]) {
	if item == nil {
		return
	}

	now := time.Now().UTC()

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.first == nil || l.first == item {
		return
	}

	if item == l.last {
		l.last = item.prev
	}

	if item.prev != nil {
		item.prev.next = item.next
	}

	if item.next != nil {
		item.next.prev = item.prev
	}

	item.prev = nil
	item.next = l.first
	l.first.prev = item
	l.first = item

	item.UpdateAt = now
}

// MoveBack moves the given node to the back of the list.
// If item is nil, the list is empty, or item is already the last node,
// MoveBack does nothing. The node's ExpireAt field is updated to the
// current time in UTC. MoveBack is safe for concurrent use.
func (l *LinkedList[T]) MoveBack(item *LinkedListNode[T]) {
	if item == nil {
		return
	}

	now := time.Now().UTC()

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.last == nil || l.last == item {
		return
	}

	if item == l.first {
		l.first = item.next
	}

	if item.prev != nil {
		item.prev.next = item.next
	}

	if item.next != nil {
		item.next.prev = item.prev
	}

	item.next = nil
	item.prev = l.last
	l.last.next = item
	l.last = item

	item.UpdateAt = now
}

// Remove deletes the given node from the list.
// If item is nil or the list is empty, Remove does nothing.
// After removal, item is detached from the list and its Next and Prev
// fields are set to nil. Remove is safe for concurrent use when other
// operations access the list through its methods.
func (l *LinkedList[T]) Remove(item *LinkedListNode[T]) {
	if item == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.first == nil {
		return
	}

	if item == l.first {
		l.first = item.next
	} else if item.prev != nil {
		item.prev.next = item.next
	}

	if item == l.last {
		l.last = item.prev
	} else if item.next != nil {
		item.next.prev = item.prev
	}

	if l.len > 0 {
		l.len--
	}

	item.next = nil
	item.prev = nil
}

// Iter calls action for each node in the list from first to last.
// If action is nil, Iter does nothing.
// The set and order of nodes are captured under a read lock before
// action is invoked, so modifications to the list inside action do
// not affect which nodes are visited.
func (l *LinkedList[T]) Iter(action func(item *LinkedListNode[T])) {
	if action == nil {
		return
	}

	l.mu.RLock()

	var nodes []*LinkedListNode[T]
	for n := l.first; n != nil; n = n.next {
		nodes = append(nodes, n)
	}

	l.mu.RUnlock()

	for _, n := range nodes {
		action(n)
	}
}

// IterFromLast calls action for each node in the list from last to first.
// If action is nil, IterFromLast does nothing.
// The set and order of nodes are captured under a read lock before
// action is invoked, so modifications to the list inside action do
// not affect which nodes are visited.
func (l *LinkedList[T]) IterFromLast(action func(item *LinkedListNode[T])) {
	if action == nil {
		return
	}

	l.mu.RLock()

	var nodes []*LinkedListNode[T]
	for n := l.last; n != nil; n = n.prev {
		nodes = append(nodes, n)
	}

	l.mu.RUnlock()

	for _, n := range nodes {
		action(n)
	}
}
