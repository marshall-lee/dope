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

var ErrWriteClosed = errors.New("buffers: write to buffer")
var ErrCloseClosed = errors.New("buffers: close of a closed buffer")

// BlockingBounded is a blocking byte buffer of a fixed capacity.
// Internally, it's implemented as a synchronized ring buffer represented
// by a pre-allocated byte slice with two offsets and is designed specifically
// for a producer-consumer use case.
//
// When the buffer is in an open state and empty, all read methods block until
// it either becomes non-empty or closed.
// When the queue is in a closed state, all read methods return immediately i.e.
// they are are non-blocking.
//
// And when the queue is in an open state and full, all write methods block
// until it becomes either partially full or empty.
// Writing to a closed queue results in panic.
type BlockingBounded struct {
	q internal.BlockingBounded[byte]
}

// NewBlockingBounded makes a new blocking byte buffer with a given capacity.
// Capacity is expected to be a positive non-zero integer value, otherwise
// this function panics.
func NewBlockingBounded(cap int) *BlockingBounded {
	if cap <= 0 {
		panic(fmt.Errorf("buffers: capacity must be a positive integer value but %v is given", cap))
	}

	var buf BlockingBounded
	buf.q.Init(cap)
	return &buf
}

// Write implements an [io.Writer] interface.
// It attempts to write the data to the buffer and it eventually
// will write every byte of it if only the buffer was not closed in the process.
// So this method blocks if we cannot write all the data at once. For example,
// if the buffer is full, this method blocks. If the buffer is only partially
// full but data length exceeds the buffer free space, this method blocks.
// And of course this method blocks if data length exceeds the buffer's fixed
// capacity.
//
// Writing to a closed buffer results in [ErrWriteClosed] error.
func (buf *BlockingBounded) Write(data []byte) (n int, err error) {
	for len(data) != 0 {
		m := buf.q.PushSomeNonEmpty(data)
		if m == 0 {
			return n, ErrWriteClosed
		}
		n += m
		data = data[m:]
	}

	// No push was really attempted so lets check that queue is closed manually.
	if n == 0 && buf.q.Closed() {
		return 0, ErrWriteClosed
	}
	return n, nil
}

// WriteByte implements an [io.ByteWriter] interface.
// It writes a single byte to the buffer.
// While the buffer is full, this methods blocks.
//
// Writing to a closed buffer results in [ErrWriteClosed] error.
func (buf *BlockingBounded) WriteByte(c byte) (err error) {
	ok := buf.q.Push(c)
	if !ok {
		err = ErrWriteClosed
	}
	return err
}

// Read implements an [io.Reader] interface.
// This methods blocks until the buffer is non-empty and fills the data slice
// with what is available at the moment.
//
// Reading from a closed empty buffer results in [io.EOF] error.
func (buf *BlockingBounded) Read(data []byte) (n int, err error) {
	if len(data) == 0 {
		if closed := buf.q.Closed(); closed {
			return 0, io.EOF
		}
		return 0, nil
	}
	n = buf.q.PopSomeNonEmpty(data)
	if n == 0 {
		err = io.EOF
	}
	return n, err
}

// ReadByte implements an [io.ByteReader] interface.
// It reads a single bytes from the buffer. While the buffer is empty,
// this method blocks.
//
// Reading from a closed empty buffer results in [io.EOF] error.
func (buf *BlockingBounded) ReadByte() (b byte, err error) {
	b, ok := buf.q.Pop()
	if !ok {
		err = io.EOF
	}
	return b, err
}

// Close puts the buffer into a closed state. After closing the buffer
// all read methods on it become non-blocking and all write methods
// on it will return [ErrWriteClosed] error.
// Calling this method twice on the same buffer results in [ErrCloseClosed] error.
func (buf *BlockingBounded) Close() error {
	if !buf.q.Close() {
		return ErrCloseClosed
	}
	return nil
}
