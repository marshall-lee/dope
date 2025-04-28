// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buffers

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type BlockingBoundedTestSuite struct {
	suite.Suite
	mu sync.RWMutex
}

func (s *BlockingBoundedTestSuite) TestNonPositiveCapPanic() {
	s.Require().PanicsWithError("buffers: capacity must be a positive integer value but 0 is given", func() { NewBlockingBounded(0) })
	s.Require().PanicsWithError("buffers: capacity must be a positive integer value but -1 is given", func() { NewBlockingBounded(-1) })
}

func (s *BlockingBoundedTestSuite) TestWriteRead() {
	buf := NewBlockingBounded(4)

	n, err := buf.Write([]byte{42, 43})
	s.Require().Equal(2, n)
	s.Require().NoError(err)

	data := make([]byte, 2)
	n, err = buf.Read(data)
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{42, 43}, data)
}

func (s *BlockingBoundedTestSuite) TestWriteFullRead() {
	buf := NewBlockingBounded(2)

	n, err := buf.Write([]byte{42, 43})
	s.Require().Equal(2, n)
	s.Require().NoError(err)

	var (
		returned bool
		nWrite   int
		writeErr error
	)
	go func() {
		nWrite, writeErr = buf.Write([]byte{44})
		s.mu.Lock()
		returned = true
		s.mu.Unlock()
	}()

	s.Require().Never(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Write unexpectedly returned")

	data := make([]byte, 1)
	n, err = buf.Read(data)
	s.Require().Equal(1, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{42}, data)

	s.Require().Eventually(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Write did not return as was expected")
	s.Require().Equal(1, nWrite)
	s.Require().NoError(writeErr)

	data = make([]byte, 2)
	n, err = buf.Read(data)
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{43, 44}, data)
}

func (s *BlockingBoundedTestSuite) TestEmptyReadAndClose() {
	buf := NewBlockingBounded(1)

	var (
		nRead    int
		readErr  error
		returned bool
	)
	go func() {
		data := make([]byte, 1)
		nRead, readErr = buf.Read(data)
		s.mu.Lock()
		returned = true
		s.mu.Unlock()
	}()

	s.Require().Never(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Read unexpectedly returned")

	err := buf.Close()
	s.Require().NoError(err)

	s.Require().Eventually(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Read did not return as was expected")
	s.Require().Zero(nRead)
	s.Require().ErrorIs(readErr, io.EOF)
}

func (s *BlockingBoundedTestSuite) TestWriteReadAndClose() {
	buf := NewBlockingBounded(3)

	var (
		returned bool
		allData  []byte
		readErr  error
	)
	go func() {
		data := make([]byte, 2)
		for {
			var n int
			if n, readErr = buf.Read(data); readErr != nil {
				break
			}
			allData = append(allData, data[:n]...)
		}
		s.mu.Lock()
		returned = true
		s.mu.Unlock()
	}()

	n, err := buf.Write([]byte{41, 42})
	s.Require().Equal(2, n)
	s.Require().NoError(err)

	s.Require().Never(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Reader loop finished unexpectedly")

	n, err = buf.Write([]byte{43})
	s.Require().Equal(1, n)
	s.Require().NoError(err)

	err = buf.Close()
	s.Require().NoError(err)

	s.Require().Eventually(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Reader loop did not finish as was expected")
	s.Require().ErrorIs(readErr, io.EOF)
	s.Require().Equal([]byte{41, 42, 43}, allData)
}

func (s *BlockingBoundedTestSuite) TestWriteFullAndReadTwice() {
	buf := NewBlockingBounded(3)

	n, err := buf.Write([]byte{42, 43, 44})
	s.Require().Equal(3, n)
	s.Require().NoError(err)

	data := make([]byte, 1)
	n, err = buf.Read(data)
	s.Require().Equal(1, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{42}, data)

	data = make([]byte, 2)
	n, err = buf.Read(data)
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{43, 44}, data)
}

func (s *BlockingBoundedTestSuite) TestWriteTwiceAndReadFull() {
	buf := NewBlockingBounded(3)

	n, err := buf.Write([]byte{42, 43})
	s.Require().Equal(2, n)
	s.Require().NoError(err)

	n, err = buf.Write([]byte{44})
	s.Require().Equal(1, n)
	s.Require().NoError(err)

	data := make([]byte, 3)
	n, err = buf.Read(data)
	s.Require().Equal(3, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{42, 43, 44}, data)
}

func (s *BlockingBoundedTestSuite) TestCloseOfAClosed() {
	buf := NewBlockingBounded(1)
	err := buf.Close()
	s.Require().NoError(err)
	err = buf.Close()
	s.Require().ErrorIs(err, ErrCloseClosed)
}

func (s *BlockingBoundedTestSuite) TestWriteToAClosed() {
	buf := NewBlockingBounded(1)
	err := buf.Close()
	s.Require().NoError(err)
	n, err := buf.Write([]byte{41, 42})
	s.Require().Equal(0, n)
	s.Require().ErrorIs(err, ErrWriteClosed)
}

func (s *BlockingBoundedTestSuite) TestBlockingWriteToAClosed() {
	buf := NewBlockingBounded(2)
	n, err := buf.Write([]byte{42})
	s.Require().Equal(1, n)
	s.Require().NoError(err)

	var (
		returned bool
		nWrite   int
		writeErr error
	)

	go func() {
		nWrite, writeErr = buf.Write([]byte{43, 44})
		s.mu.Lock()
		returned = true
		s.mu.Unlock()
	}()

	s.Require().Never(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Write returned unexpectedly")

	err = buf.Close()
	s.Require().NoError(err)

	s.Require().Eventually(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Write did not return as was expected")
	s.Require().Equal(1, nWrite)
	s.Require().ErrorIs(writeErr, ErrWriteClosed)
}

func (s *BlockingBoundedTestSuite) TestInterface() {
	s.Require().Implements((*Interface)(nil), NewBlockingBounded(1))
}

func TestBlockingBounded(t *testing.T) {
	suite.Run(t, new(BlockingBoundedTestSuite))
}
