package structure_test

import (
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/platform/structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinkedList_Empty(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	require.Equal(t, 0, list.Len(), "list length should be 0")

	forward := collectForward(list)
	assert.Empty(t, forward, "expected no items in forward iteration")

	backward := collectBackward(list)
	assert.Empty(t, backward, "expected no items in backward iteration")
}

func TestPushFront_OrderAndLen(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	list.PushFront(1)
	list.PushFront(2)
	list.PushFront(3)

	require.Equal(t, 3, list.Len(), "list length should be 3")

	forward := collectForward(list)
	assert.Equal(t, []int{3, 2, 1}, forward, "forward iteration mismatch")

	backward := collectBackward(list)
	assert.Equal(t, []int{1, 2, 3}, backward, "backward iteration mismatch")
}

func TestPushBack_OrderAndLen(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	require.Equal(t, 3, list.Len(), "list length should be 3")

	forward := collectForward(list)
	assert.Equal(t, []int{1, 2, 3}, forward, "forward iteration mismatch")

	backward := collectBackward(list)
	assert.Equal(t, []int{3, 2, 1}, backward, "backward iteration mismatch")
}

func TestMoveFront_MiddleNode(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	list.PushBack(1)
	node2 := list.PushBack(2)
	list.PushBack(3)

	list.MoveFront(node2)

	require.Equal(t, 3, list.Len(), "list length should be 3")

	forward := collectForward(list)
	assert.Equal(t, []int{2, 1, 3}, forward, "forward iteration mismatch")

	backward := collectBackward(list)
	assert.Equal(t, []int{3, 1, 2}, backward, "backward iteration mismatch")
}

func TestMoveBack_MiddleNode(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	list.PushBack(1)
	node2 := list.PushBack(2)
	list.PushBack(3)

	list.MoveBack(node2)

	require.Equal(t, 3, list.Len(), "list length should be 3")

	forward := collectForward(list)
	assert.Equal(t, []int{1, 3, 2}, forward, "forward iteration mismatch")

	backward := collectBackward(list)
	assert.Equal(t, []int{2, 3, 1}, backward, "backward iteration mismatch")
}

func TestMoveFront_FirstNode_NoChangeOrderAndBothDirectionsConsistent(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	node1 := list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	list.MoveFront(node1)

	require.Equal(t, 3, list.Len(), "list length should be 3")

	forward := collectForward(list)
	assert.Equal(t, []int{1, 2, 3}, forward, "forward iteration mismatch")

	backward := collectBackward(list)
	assert.Equal(t, []int{3, 2, 1}, backward, "backward iteration mismatch")
}

func TestMoveBack_LastNode_NoChangeOrderAndBothDirectionsConsistent(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	list.PushBack(1)
	list.PushBack(2)
	node3 := list.PushBack(3)

	list.MoveBack(node3)

	require.Equal(t, 3, list.Len(), "list length should be 3")

	forward := collectForward(list)
	assert.Equal(t, []int{1, 2, 3}, forward, "forward iteration mismatch")

	backward := collectBackward(list)
	assert.Equal(t, []int{3, 2, 1}, backward, "backward iteration mismatch")
}

func TestRemove_MiddleNode(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	list.PushBack(1)
	node2 := list.PushBack(2)
	list.PushBack(3)

	list.Remove(node2)

	require.Equal(t, 2, list.Len(), "list length should be 2 after remove")

	forward := collectForward(list)
	assert.Equal(t, []int{1, 3}, forward, "forward iteration mismatch")

	backward := collectBackward(list)
	assert.Equal(t, []int{3, 1}, backward, "backward iteration mismatch")
}

func TestRemove_FirstAndLastNodes(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	node1 := list.PushBack(1)
	list.PushBack(2)
	node3 := list.PushBack(3)

	list.Remove(node1)

	require.Equal(t, 2, list.Len(), "list length should be 2 after removing first")

	forward := collectForward(list)
	assert.Equal(t, []int{2, 3}, forward, "forward iteration after removing first mismatch")

	backward := collectBackward(list)
	assert.Equal(t, []int{3, 2}, backward, "backward iteration after removing first mismatch")

	list.Remove(node3)

	require.Equal(t, 1, list.Len(), "list length should be 1 after removing last")

	forward = collectForward(list)
	assert.Equal(t, []int{2}, forward, "forward iteration after removing last mismatch")

	backward = collectBackward(list)
	assert.Equal(t, []int{2}, backward, "backward iteration after removing last mismatch")
}

func TestRemove_SingleElementList(t *testing.T) {
	t.Parallel()

	list := structure.NewLinkedList[int]()

	node := list.PushBack(42)

	require.Equal(t, 1, list.Len(), "list length should be 1")

	list.Remove(node)

	require.Equal(t, 0, list.Len(), "list length should be 0 after remove")

	forward := collectForward(list)
	assert.Empty(t, forward, "expected no items in forward iteration")

	backward := collectBackward(list)
	assert.Empty(t, backward, "expected no items in backward iteration")
}

func collectForward[T any](l *structure.LinkedList[T]) []T {
	var res []T

	l.Iter(func(node *structure.LinkedListNode[T]) {
		res = append(res, node.Item)
	})

	return res
}

func collectBackward[T any](l *structure.LinkedList[T]) []T {
	var res []T

	l.IterFromLast(func(node *structure.LinkedListNode[T]) {
		res = append(res, node.Item)
	})

	return res
}
