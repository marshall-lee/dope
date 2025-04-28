// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"fmt"

	internal "github.com/marshall-lee/dope/internal/queues"
)

// Bounded is a FIFO queue of a fixed capacity.
// Internally, it's implemented as a simple ring buffer represented by a
// pre-allocated slice with two offsets.
// All methods are non-blocking and this queue is not suitable for usage by multiple goroutines.
type Bounded[T any] struct {
	q internal.Bounded[T]
}

// NewBounded makes a new FIFO queue with a given capacity.
// The buffer is allocated here and is never reallocated.
//
// Capacity is expected to be a positive non-zero integer value, otherwise
// this function panics.
func NewBounded[T any](cap int) *Bounded[T] {
	if cap <= 0 {
		panic(fmt.Errorf("queues: capacity must be a positive integer value but %v is given", cap))
	}

	var queue Bounded[T]
	queue.q.Init(cap)
	return &queue
}

// Push attemtps to add an element to the queue. Return value is true
// if there was a room to add a new element. If the queue is full,
// then return value is false.
func (queue *Bounded[T]) Push(elem T) (ok bool) {
	return queue.q.Push(elem)
}

// PushSome attemtps to add to the queue at least some of the slice elements.
// Return value is the number of slice elements added to the queue. Zero means
// that the queue is full.
func (queue *Bounded[T]) PushSome(values []T) (n int) {
	return queue.q.PushSome(values)
}

// Pop attempts to consume one element from the queue. Return value ok is true
// whenever the element is successfully consumed. Otherwise, return value ok is
// false and it basically means that the queue is empty.
func (queue *Bounded[T]) Pop() (value T, ok bool) {
	return queue.q.Pop()
}

// PopSome attempts to consume a bunch of elements from the queue.
// Return value is the number of elements consumed from it.
func (queue *Bounded[T]) PopSome(out []T) (n int) {
	return queue.q.PopSome(out)
}

// Cap returns the capacity of the queue.
func (queue *Bounded[T]) Cap() int {
	return queue.q.Cap()
}

// Len returns the current number of elements in the queue.
func (queue *Bounded[T]) Len() int {
	return queue.q.Len()
}

// Available returns how much elements the queue is able to take.
func (queue *Bounded[T]) Available() int {
	return queue.q.Available()
}

// Full returns true if the queue is full.
func (queue *Bounded[T]) Full() bool {
	return queue.q.Full()
}

// Empty returns true if the queue is empty.
func (queue *Bounded[T]) Empty() bool {
	return queue.q.Empty()
}

// Slice returns a copy of the queue contents.
func (queue *Bounded[T]) Slice() []T {
	return queue.q.Slice()
}
