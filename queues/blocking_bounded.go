// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"fmt"

	internal "github.com/marshall-lee/dope/internal/queues"
)

// BlockingBounded is a blocking FIFO queue of a fixed capacity.
// Internally, it's implemented as a synchronized ring buffer represented
// by a pre-allocated slice with two offsets and is designed specifically
// for a producer-consumer use case.
//
// When the queue is in an open state and empty, all read methods block until
// it either becomes non-empty or closed.
// When the queue is in a closed state, all read methods return immediately i.e.
// they are are non-blocking.
//
// And when the queue is in an open state and full, all write methods block
// until it becomes either partially full or empty.
// Writing to a closed queue results in panic.
type BlockingBounded[T any] struct {
	q internal.BlockingBounded[T]
}

// NewBlockingBounded makes a new blocking FIFO queue with a given capacity.
// Capacity is expected to be a positive non-zero integer value, otherwise
// this function panics.
func NewBlockingBounded[T any](cap int) *BlockingBounded[T] {
	if cap <= 0 {
		panic(fmt.Errorf("queues: capacity must be a positive integer value but %v is given", cap))
	}

	var queue BlockingBounded[T]
	queue.q.Init(cap)
	return &queue
}

// Push adds an element to the queue. This method blocks while the queue
// is full and panics if the queue is closed.
func (queue *BlockingBounded[T]) Push(elem T) {
	if ok := queue.q.Push(elem); !ok {
		panic(ErrPushClosed)
	}
}

// PushSome adds to the queue at least some of the slice elements. This method
// blocks while the queue is full and panics if the queue is closed. Return value is the number of slice
// elements added to the queue and it's always non-zero.
func (queue *BlockingBounded[T]) PushSome(values []T) (n int) {
	if len(values) == 0 {
		if queue.q.Closed() {
			panic(ErrPushClosed)
		}
		return 0
	}
	if n = queue.q.PushSomeNonEmpty(values); n == 0 {
		panic(ErrPushClosed)
	}
	return n
}

// PushAll adds to the queue all the elements from the slice. This method blocks until
// eventually the room is found for every element of the slice. If in the process of pushing elements
// to the queue it was closed, this method panics.
func (queue *BlockingBounded[T]) PushAll(values []T) {
	var pushed bool
	for len(values) != 0 {
		m := queue.q.PushSomeNonEmpty(values)
		if m == 0 {
			panic(ErrPushClosed)
		}
		values = values[m:]
		pushed = true
	}

	// No push was really attempted so lets check that queue is closed manually.
	if !pushed && queue.q.Closed() {
		panic(ErrPushClosed)
	}
}

// Pop attempts to consume one element from the queue. This method blocks while the queue is empty
// and returns immediately if the queue is closed.
func (queue *BlockingBounded[T]) Pop() (value T, ok bool) {
	return queue.q.Pop()
}

// PopSome attempts to consume a bunch of elements from the queue. This method blocks while the
// queue is empty and returns immediately if the queue is closed. Return value is the number of
// elements consumed from the queue.
func (queue *BlockingBounded[T]) PopSome(out []T) (n int) {
	if len(out) == 0 {
		return 0
	}
	return queue.q.PopSomeNonEmpty(out)
}

// Close puts the queue into a closed state. After closing the queue
// all read methods on it become non-blocking and all write methods
// on it will panic.
// Calling this method twice on the same queue also results in panic.
func (queue *BlockingBounded[T]) Close() {
	if ok := queue.q.Close(); !ok {
		panic(ErrCloseClosed)
	}
}
