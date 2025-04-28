// Copyright 2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

type Unbounded[T any] struct {
	elems []T
	offset int
}

func (queue *Unbounded[T]) Init() {}

func (queue *Unbounded[T]) Push(elem T) {
	panic("TODO")
}

func (queue *Unbounded[T]) PushAll(values []T) {
	panic("TODO")
}

func (queue *Unbounded[T]) Pop() (value T, ok bool) {
	panic("TODO")
}

func (queue *Unbounded[T]) PopSome(out []T) (n int) {
	panic("TODO")
}

func (queue *Unbounded[T]) Len() int {
	return len(queue.elems) - queue.offset
}

func (queue *Unbounded[T]) Empty() bool {
	return len(queue.elems) == queue.offset
}
