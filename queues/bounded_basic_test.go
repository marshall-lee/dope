// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// BoundedBasicTestSuite tests Push/Pop methods.
type BoundedBasicTestSuite struct {
	suite.Suite
}

func (s *BoundedBasicTestSuite) TestNegativeCapPanic() {
	s.Require().PanicsWithError("queues: capacity must be a positive integer value but 0 is given", func() { NewBounded[int](0) })
	s.Require().PanicsWithError("queues: capacity must be a positive integer value but -1 is given", func() { NewBounded[int](-1) })
}

func (s *BoundedBasicTestSuite) TestEmptyPop() {
	queue := NewBounded[int](1)

	s.Require().True(queue.Empty())
	_, ok := queue.Pop()
	s.Require().False(ok)
}

func (s *BoundedBasicTestSuite) TestPush() {
	queue := NewBounded[int](2)

	ok := queue.Push(42)
	s.Require().True(ok)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	ok = queue.Push(43)
	s.Require().True(ok)
	s.Require().False(queue.Empty())
	s.Require().True(queue.Full())

	// Overflow
	ok = queue.Push(44)
	s.Require().False(ok)
}

func (s *BoundedBasicTestSuite) TestMultiplePushAndPop() {
	queue := NewBounded[int](3)

	ok := queue.Push(42)
	s.Require().True(ok)

	ok = queue.Push(43)
	s.Require().True(ok)

	value, ok := queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(42, value)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	ok = queue.Push(44)
	s.Require().True(ok)
	s.Require().False(queue.Full())

	ok = queue.Push(45)
	s.Require().True(ok)
	s.Require().True(queue.Full())

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(43, value)
	s.Require().False(queue.Empty())

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(44, value)
	s.Require().False(queue.Empty())

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(45, value)
	s.Require().True(queue.Empty())
	s.Require().False(queue.Full())
}

func (s *BoundedBasicTestSuite) TestMultiplePushAndPop2() {
	queue := NewBounded[int](3)

	ok := queue.Push(42)
	s.Require().True(ok)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	ok = queue.Push(43)
	s.Require().True(ok)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	ok = queue.Push(44)
	s.Require().True(ok)
	s.Require().False(queue.Empty())
	s.Require().True(queue.Full())

	value, ok := queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(42, value)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(43, value)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	ok = queue.Push(45)
	s.Require().True(ok)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	ok = queue.Push(46)
	s.Require().True(ok)
	s.Require().False(queue.Empty())
	s.Require().True(queue.Full())

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(44, value)
	s.Require().False(queue.Empty())

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(45, value)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	value, ok = queue.Pop()
	s.Require().True(ok)
	s.Require().Equal(46, value)
	s.Require().True(queue.Empty())
	s.Require().False(queue.Full())
}

func TestBoundedBasic(t *testing.T) {
	suite.Run(t, new(BoundedBasicTestSuite))
}
