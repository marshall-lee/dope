// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buffers

import (
	"errors"
	"fmt"
	"io"

	internal "github.com/marshall-lee/dope/internal/queues"
)

// Bounded is a byte buffer of a fixed capacity.
// Internally, it's implemented as a simple ring buffer represented by a
// pre-allocated byte slice with two offsets.
// All methods are non-blocking and this buffer is not suitable for usage by
// multiple goroutines.
type Bounded struct {
	q internal.Bounded[byte]
}

var ErrBufferIsFull = errors.New("buffers: buffer is full")

// NewBounded makes a new bytes buffer with a given capacity.
// The buffer is allocated here and is never reallocated.
//
// Capacity is expected to be a positive non-zero integer value, otherwise
// this function panics.
func NewBounded(cap int) *Bounded {
	if cap <= 0 {
		panic(fmt.Errorf("buffers: capacity must be a positive integer value but %v is given", cap))
	}

	var buf Bounded
	buf.q.Init(cap)
	return &buf
}

// Write implements an [io.Writer] interface.
// If data length exceeds the [Bounded.Available], then the data
// is written only partially, and the [ErrBufferIsFull] is returned.
func (buf *Bounded) Write(data []byte) (n int, err error) {
	n = buf.q.PushSome(data)
	if n < len(data) {
		err = ErrBufferIsFull
	}
	return n, err
}

// WriteByte implements an [io.ByteWriter] interface.
// It writes a single byte to the buffer.
// If the buffer is full, this method returns [ErrBufferIsFull].
func (buf *Bounded) WriteByte(c byte) (err error) {
	ok := buf.q.Push(c)
	if !ok {
		err = ErrBufferIsFull
	}
	return err
}

// Read implements an [io.Reader] interface.
// It fills the data slice with what is available at the moment.
// If the buffer is empty, this method returns [io.EOF].
func (buf *Bounded) Read(data []byte) (n int, err error) {
	n = buf.q.PopSome(data)
	if n == 0 && (len(data) != 0 || buf.q.Empty()) {
		err = io.EOF
	}
	return n, err
}

// ReadByte implements an [io.ByteReader] interface.
// It reads a single byte from the buffer.
// If the buffer is empty, this method returns [io.EOF].
func (buf *Bounded) ReadByte() (b byte, err error) {
	b, ok := buf.q.Pop()
	if !ok {
		err = io.EOF
	}
	return b, err
}

// Cap returns the capacity of the buffer.
func (buf *Bounded) Cap() int {
	return buf.q.Cap()
}

// Len returns the size of unread data in the buffer.
func (buf *Bounded) Len() int {
	return buf.q.Len()
}

// Available returns buffer's free space i.e. how much bytes
// is available to write to the buffer.
func (buf *Bounded) Available() int {
	return buf.q.Available()
}

// Full returns true if the buffer is full.
func (buf *Bounded) Full() bool {
	return buf.q.Full()
}

// Full returns true if the buffer is empty.
func (buf *Bounded) Empty() bool {
	return buf.q.Empty()
}

// Bytes returns a copy of the unread data in the buffer.
func (buf *Bounded) Bytes() []byte {
	return buf.q.Slice()
}
