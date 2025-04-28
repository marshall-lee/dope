// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// BoundedBatchTestSuite tests batch methods PushSome/PopSome in different cases.
type BoundedBatchTestSuite struct {
	suite.Suite
}

func (s *BoundedBatchTestSuite) TestEmptyPop() {
	queue := NewBounded[int](1)

	s.Require().True(queue.Empty())
	data := make([]int, 1)
	n := queue.PopSome(data)
	s.Require().Zero(n)
}

func (s *BoundedBatchTestSuite) TestPush() {
	queue := NewBounded[int](2)

	n := queue.PushSome([]int{42})
	s.Require().Equal(1, n)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())
}

func (s *BoundedBatchTestSuite) TestPushSomeFull() {
	queue := NewBounded[int](2)

	n := queue.PushSome([]int{42, 43})
	s.Require().Equal(2, n)
	s.Require().False(queue.Empty())
	s.Require().True(queue.Full())
}

func (s *BoundedBatchTestSuite) TestPushSomeOverflow() {
	queue := NewBounded[int](2)

	n := queue.PushSome([]int{42, 43, 44})
	s.Require().Equal(2, n)
	s.Require().True(queue.Full())
}

func (s *BoundedBatchTestSuite) TestPushSomeFullAndPopTwice() {
	queue := NewBounded[int](3)

	n := queue.PushSome([]int{42, 43, 44})
	s.Require().Equal(3, n)
	s.Require().True(queue.Full())

	data := make([]int, 1)
	n = queue.PopSome(data)
	s.Require().Equal(1, n)
	s.Require().Equal([]int{42}, data)
	s.Require().False(queue.Full())

	data = make([]int, 2)
	n = queue.PopSome(data)
	s.Require().Equal(2, n)
	s.Require().Equal([]int{43, 44}, data)
	s.Require().False(queue.Full())
}

func (s *BoundedBatchTestSuite) TestPushSomeTwiceAndPopFull() {
	queue := NewBounded[int](3)

	n := queue.PushSome([]int{42, 43})
	s.Require().Equal(2, n)
	s.Require().False(queue.Full())

	n = queue.PushSome([]int{44})
	s.Require().Equal(1, n)
	s.Require().True(queue.Full())

	data := make([]int, 3)
	n = queue.PopSome(data)
	s.Require().Equal(3, n)
	s.Require().Equal([]int{42, 43, 44}, data)
	s.Require().False(queue.Full())
}

func (s *BoundedBatchTestSuite) TestPushSomeAndPop() {
	queue := NewBounded[int](3)
	data := make([]int, 3)

	n := queue.PushSome([]int{42, 43})
	s.Require().Equal(2, n)
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	n = queue.PopSome(data[:1])
	s.Require().Equal(1, n)
	s.Require().Equal([]int{42}, data[:n])
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	n = queue.PushSome([]int{44, 45})
	s.Require().Equal(2, n)
	s.Require().False(queue.Empty())
	s.Require().True(queue.Full())

	n = queue.PopSome(data)
	s.Require().Equal(3, n)
	s.Require().Equal([]int{43, 44, 45}, data[:n])
	s.Require().True(queue.Empty())
	s.Require().False(queue.Full())
}

func (s *BoundedBatchTestSuite) TestPushSomeAndPop2() {
	queue := NewBounded[int](3)
	data := make([]int, 3)

	n := queue.PushSome([]int{42, 43, 44})
	s.Require().Equal(3, n)
	s.Require().False(queue.Empty())
	s.Require().True(queue.Full())

	n = queue.PopSome(data[:2])
	s.Require().Equal(2, n)
	s.Require().Equal([]int{42, 43}, data[:n])
	s.Require().False(queue.Full())

	n = queue.PushSome([]int{45, 46})
	s.Require().Equal(2, n)
	s.Require().False(queue.Empty())
	s.Require().True(queue.Full())

	n = queue.PopSome(data[:2])
	s.Require().Equal(2, n)
	s.Require().Equal([]int{44, 45}, data[:n])
	s.Require().False(queue.Empty())
	s.Require().False(queue.Full())

	n = queue.PopSome(data)
	s.Require().Equal(1, n)
	s.Require().Equal([]int{46}, data[:n])
	s.Require().True(queue.Empty())
	s.Require().False(queue.Full())
}

func TestBoundedBatch(t *testing.T) {
	suite.Run(t, new(BoundedBatchTestSuite))
}
