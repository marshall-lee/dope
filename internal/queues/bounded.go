// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

type Bounded[T any] struct {
	elems   []T
	wOffset int
	rOffset int
	full    bool
}

func (queue *Bounded[T]) Init(cap int) {
	queue.elems = make([]T, cap, cap)
}

func (queue *Bounded[T]) Push(elem T) (ok bool) {
	if queue.full {
		return false
	}
	queue.elems[queue.wOffset] = elem
	if queue.wOffset += 1; queue.wOffset == len(queue.elems) {
		queue.wOffset = 0
	}
	if queue.wOffset == queue.rOffset {
		queue.full = true
	}
	return true
}

func (queue *Bounded[T]) PushSome(values []T) (n int) {
	if queue.wOffset >= queue.rOffset {
		if queue.full {
			return 0
		}
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
	if n > 0 && queue.wOffset == queue.rOffset {
		queue.full = true
	}
	return n
}

func (queue *Bounded[T]) Pop() (value T, ok bool) {
	if queue.wOffset == queue.rOffset && !queue.full {
		return value, false
	}
	value = queue.elems[queue.rOffset]
	if queue.rOffset += 1; queue.rOffset == len(queue.elems) {
		queue.rOffset = 0
	}
	queue.full = false
	return value, true
}

func (queue *Bounded[T]) PopSome(out []T) (n int) {
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
	if n > 0 {
		queue.full = false
	}
	return n
}

func (queue *Bounded[T]) Slice() []T {
	var out []T

	if queue.wOffset < queue.rOffset || queue.full {
		out = make([]T, len(queue.elems)-(queue.rOffset-queue.wOffset))
		n := copy(out, queue.elems[queue.rOffset:])
		copy(out[n:], queue.elems[:queue.wOffset])
	} else if queue.wOffset > queue.rOffset {
		out = make([]T, queue.wOffset-queue.rOffset)
		copy(out, queue.elems[queue.rOffset:queue.wOffset])
	}

	return out
}

func (queue *Bounded[T]) Cap() int {
	return len(queue.elems)
}

func (queue *Bounded[T]) Len() int {
	if queue.full {
		return len(queue.elems)
	} else if queue.wOffset >= queue.rOffset {
		return queue.wOffset - queue.rOffset
	} else {
		return len(queue.elems) - (queue.rOffset - queue.wOffset)
	}
}

func (queue *Bounded[T]) Available() int {
	if queue.full {
		return 0
	} else if queue.wOffset >= queue.rOffset {
		return len(queue.elems) - (queue.wOffset - queue.rOffset)
	} else {
		return queue.rOffset - queue.wOffset
	}
}

func (queue *Bounded[T]) Full() bool {
	return queue.full
}

func (queue *Bounded[T]) Empty() bool {
	return !queue.full && queue.wOffset == queue.rOffset
}
