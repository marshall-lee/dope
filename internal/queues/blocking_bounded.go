// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"sync"
)

type BlockingBounded[T any] struct {
	mu      sync.RWMutex
	rCond   *sync.Cond
	wCond   *sync.Cond
	elems   []T
	wOffset int
	rOffset int
	full    bool
	closed  bool
}

func (queue *BlockingBounded[T]) Init(cap int) {
	queue.elems = make([]T, cap, cap)
	queue.rCond = sync.NewCond(&queue.mu)
	queue.wCond = sync.NewCond(&queue.mu)
}

func (queue *BlockingBounded[T]) Push(elem T) (ok bool) {
	queue.mu.Lock()
	if !queue.waitWriteable() {
		queue.mu.Unlock()
		return false
	}
	queue.elems[queue.wOffset] = elem
	if queue.wOffset += 1; queue.wOffset == len(queue.elems) {
		queue.wOffset = 0
	}
	if queue.wOffset == queue.rOffset {
		queue.full = true
	} else {
		queue.wCond.Signal()
	}
	queue.rCond.Signal()
	queue.mu.Unlock()
	return true
}

func (queue *BlockingBounded[T]) PushSomeNonEmpty(values []T) (n int) {
	queue.mu.Lock()
	if !queue.waitWriteable() {
		queue.mu.Unlock()
		return 0
	}
	if queue.wOffset >= queue.rOffset {
		n = copy(queue.elems[queue.wOffset:], values)
		if n < len(values) {
			m := copy(queue.elems[:queue.rOffset], values[n:])
			n += m
			queue.wOffset = m
		} else {
			if queue.wOffset += n; queue.wOffset == len(queue.elems) {
				queue.wOffset = 0
			}
		}
	} else {
		n = copy(queue.elems[queue.wOffset:queue.rOffset], values)
		queue.wOffset += n
	}
	if queue.wOffset == queue.rOffset {
		queue.full = true
	} else {
		queue.wCond.Signal()
	}
	queue.rCond.Signal()
	queue.mu.Unlock()
	return n
}

func (queue *BlockingBounded[T]) Pop() (value T, ok bool) {
	queue.mu.Lock()
	if !queue.waitReadable() {
		queue.mu.Unlock()
		return value, false
	}
	value = queue.elems[queue.rOffset]
	if queue.rOffset += 1; queue.rOffset == len(queue.elems) {
		queue.rOffset = 0
	}
	queue.full = false
	if queue.wOffset != queue.rOffset {
		queue.rCond.Signal()
	}
	queue.wCond.Signal()
	queue.mu.Unlock()
	return value, true
}

func (queue *BlockingBounded[T]) PopSomeNonEmpty(out []T) (n int) {
	queue.mu.Lock()
	if !queue.waitReadable() {
		queue.mu.Unlock()
		return 0
	}
	if queue.wOffset < queue.rOffset || queue.full {
		n = copy(out, queue.elems[queue.rOffset:])
		if n < len(out) {
			m := copy(out[n:], queue.elems[:queue.wOffset])
			n += m
			queue.rOffset = m
		} else {
			if queue.rOffset += n; queue.rOffset == len(queue.elems) {
				queue.rOffset = 0
			}
		}
	} else {
		n = copy(out, queue.elems[queue.rOffset:queue.wOffset])
		queue.rOffset += n
	}
	queue.full = false
	if queue.wOffset != queue.rOffset {
		queue.rCond.Signal()
	}
	queue.wCond.Signal()
	queue.mu.Unlock()
	return n
}

func (queue *BlockingBounded[T]) Close() (ok bool) {
	queue.mu.Lock()
	if ok = !queue.closed; ok {
		queue.closed = true
		queue.rCond.Signal()
		queue.wCond.Signal()
	}
	queue.mu.Unlock()
	return ok
}

func (queue *BlockingBounded[T]) Closed() bool {
	queue.mu.RLock()
	result := queue.closed
	queue.mu.RUnlock()
	return result
}

func (queue *BlockingBounded[T]) waitWriteable() (ok bool) {
	for queue.full {
		if queue.closed {
			return false
		}
		queue.wCond.Wait()
	}
	return !queue.closed
}

func (queue *BlockingBounded[T]) waitReadable() (ok bool) {
	for queue.empty() {
		if queue.closed {
			return false
		}
		queue.rCond.Wait()
	}
	return true
}

func (queue *BlockingBounded[T]) cap() int {
	return len(queue.elems)
}

func (queue *BlockingBounded[T]) len() int {
	if queue.full {
		return len(queue.elems)
	} else if queue.wOffset >= queue.rOffset {
		return queue.wOffset - queue.rOffset
	} else {
		return len(queue.elems) - (queue.rOffset - queue.wOffset)
	}
}

func (queue *BlockingBounded[T]) available() int {
	if queue.full {
		return 0
	} else if queue.wOffset >= queue.rOffset {
		return len(queue.elems) - (queue.wOffset - queue.rOffset)
	} else {
		return queue.rOffset - queue.wOffset
	}
}

func (queue *BlockingBounded[T]) empty() bool {
	return queue.wOffset == queue.rOffset && !queue.full
}
