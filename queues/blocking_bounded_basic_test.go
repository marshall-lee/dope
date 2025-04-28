// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"sync"
	"testing"
	"time"

	internal "github.com/marshall-lee/dope/internal/tests"
	"github.com/stretchr/testify/suite"
)

// BlockingBoundedBasicTestSuite tests Push/Pop methods.
type BlockingBoundedBasicTestSuite struct {
	suite.Suite
	mu sync.RWMutex
}

func (s *BlockingBoundedBasicTestSuite) TestNonPositiveCapPanic() {
	s.Require().PanicsWithError("queues: capacity must be a positive integer value but 0 is given", func() { NewBlockingBounded[int](0) })
	s.Require().PanicsWithError("queues: capacity must be a positive integer value but -1 is given", func() { NewBlockingBounded[int](-1) })
}


func (s *BlockingBoundedBasicTestSuite) TestPushPop() {
	queue := NewBlockingBounded[int](2)
	queue.Push(42)
	queue.Push(43)

	val, ok := queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(42, val)

	val, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(43, val)
}

func (s *BlockingBoundedBasicTestSuite) TestPushFullPop() {
	queue := NewBlockingBounded[int](2)
	queue.Push(42)
	queue.Push(43)

	var returned bool
	go func() {
		queue.Push(44)
		s.mu.Lock()
		returned = true
		s.mu.Unlock()
	}()

	s.Require().Never(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Push unexpectedly returned")

	val, ok := queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(42, val)

	s.Require().Eventually(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Push did not return as was expected")

	val, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(43, val)

	val, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(44, val)
}

func (s *BlockingBoundedBasicTestSuite) TestEmptyPopAndClose() {
	queue := NewBlockingBounded[int](1)

	var (
		popOk    bool
		returned bool
	)
	go func() {
		_, popOk = queue.Pop()
		s.mu.Lock()
		returned = true
		s.mu.Unlock()
	}()

	s.Require().Never(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Pop unexpectedly returned")

	queue.Close()

	s.Require().Eventually(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Pop did not return as was expected")
	s.Require().False(popOk)
}

func (s *BlockingBoundedBasicTestSuite) TestPushPopAndClose() {
	queue := NewBlockingBounded[int](3)

	var (
		returned bool
		values   []int
	)
	go func() {
		for {
			val, ok := queue.Pop()
			if !ok {
				break
			}
			values = append(values, val)
		}
		s.mu.Lock()
		returned = true
		s.mu.Unlock()
	}()

	queue.Push(41)
	queue.Push(42)

	s.Require().Never(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Consumer loop finished unexpectedly")

	queue.Push(43)
	queue.Close()

	s.Require().Eventually(func() bool {
		s.mu.RLock()
		ok := returned
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Consumer loop did not finish as was expected")
	s.Require().Equal([]int{41, 42, 43}, values)
}

func (s *BlockingBoundedBasicTestSuite) TestMultiplePushAndPop() {
	queue := NewBlockingBounded[int](3)
	queue.Push(42)
	queue.Push(43)

	value, ok := queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(42, value)

	queue.Push(44)
	queue.Push(45)

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(43, value)

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(44, value)

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(45, value)
}

func (s *BlockingBoundedBasicTestSuite) TestMultiplePushAndPop2() {
	queue := NewBlockingBounded[int](3)

	queue.Push(42)
	queue.Push(43)
	queue.Push(44)

	value, ok := queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(42, value)

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(43, value)

	queue.Push(45)
	queue.Push(46)

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(44, value)

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(45, value)

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(46, value)
}

func (s *BlockingBoundedBasicTestSuite) TestCloseOfAClosed() {
	queue := NewBlockingBounded[int](1)
	queue.Close()
	s.Require().PanicsWithValue(ErrCloseClosed, func() { queue.Close() }, "Close did not panic as was expected")
}

func (s *BlockingBoundedBasicTestSuite) TestPushToAClosed() {
	queue := NewBlockingBounded[int](1)
	queue.Close()
	s.Require().PanicsWithValue(ErrPushClosed, func() { queue.Push(42) }, "Push did not panic as was expected")
}

func (s *BlockingBoundedBasicTestSuite) TestBlockingPushToAClosed() {
	queue := NewBlockingBounded[int](1)
	queue.Push(42)

	var (
		panicked bool
		err      any
	)

	s.goCapturePanic(func() { queue.Push(43) }, &panicked, &err)

	s.Require().Never(func() bool {
		s.mu.RLock()
		ok := panicked
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Push panicked unexpectedly")

	queue.Close()

	s.Require().Eventually(func() bool {
		s.mu.RLock()
		ok := panicked
		s.mu.RUnlock()
		return ok
	}, 200*time.Millisecond, 10*time.Millisecond, "Push did not panic as was expected")
}

func (s *BlockingBoundedBasicTestSuite) goCapturePanic(f func(), panickedPtr *bool, errPtr *any) {
	internal.GoCapturePanic(f, &s.mu, panickedPtr, errPtr)
}

func TestBlockingBoundedBasic(t *testing.T) {
	suite.Run(t, new(BlockingBoundedBasicTestSuite))
}
