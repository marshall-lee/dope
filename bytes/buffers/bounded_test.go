package buffers

import (
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BoundedTestSuite struct {
	suite.Suite
}

func (s *BoundedTestSuite) TestEmptyRead() {
	buf := NewBounded(1)

	s.Require().True(buf.Empty())
	data := make([]byte, 1)
	n, err := buf.Read(data)
	s.Require().Zero(n)
	s.Require().ErrorIs(err, io.EOF)
}

func (s *BoundedTestSuite) TestWrite() {
	buf := NewBounded(2)

	n, err := buf.Write([]byte{42})
	s.Require().Equal(1, n)
	s.Require().NoError(err)
	s.Require().False(buf.Empty())
	s.Require().False(buf.Full())
}

func (s *BoundedTestSuite) TestWriteFull() {
	buf := NewBounded(2)

	n, err := buf.Write([]byte{42, 43})
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().False(buf.Empty())
	s.Require().True(buf.Full())
}

func (s *BoundedTestSuite) TestWriteOverflow() {
	buf := NewBounded(2)

	n, err := buf.Write([]byte{42, 43, 44})
	s.Require().Equal(2, n)
	s.Require().ErrorIs(err, ErrBufferIsFull)
	s.Require().True(buf.Full())
}

func (s *BoundedTestSuite) TestWriteFullAndReadTwice() {
	buf := NewBounded(3)

	n, err := buf.Write([]byte{42, 43, 44})
	s.Require().Equal(3, n)
	s.Require().NoError(err)
	s.Require().True(buf.Full())

	data := make([]byte, 1)
	n, err = buf.Read(data)
	s.Require().Equal(1, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{42}, data)
	s.Require().False(buf.Full())

	data = make([]byte, 2)
	n, err = buf.Read(data)
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{43, 44}, data)
	s.Require().False(buf.Full())
}

func (s *BoundedTestSuite) TestWriteTwiceAndReadFull() {
	buf := NewBounded(3)

	n, err := buf.Write([]byte{42, 43})
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().False(buf.Full())

	n, err = buf.Write([]byte{44})
	s.Require().Equal(1, n)
	s.Require().NoError(err)
	s.Require().True(buf.Full())

	data := make([]byte, 3)
	n, err = buf.Read(data)
	s.Require().Equal(3, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{42, 43, 44}, data)
	s.Require().False(buf.Full())
}

func (s *BoundedTestSuite) TestWriteAndRead() {
	buf := NewBounded(3)
	data := make([]byte, 3)

	n, err := buf.Write([]byte{42, 43})
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().False(buf.Empty())
	s.Require().False(buf.Full())

	n, err = buf.Read(data[:1])
	s.Require().Equal(1, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{42}, data[:n])
	s.Require().False(buf.Empty())
	s.Require().False(buf.Full())

	n, err = buf.Write([]byte{44, 45})
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().False(buf.Empty())
	s.Require().True(buf.Full())

	n, err = buf.Read(data)
	s.Require().Equal(3, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{43, 44, 45}, data[:n])
	s.Require().True(buf.Empty())
	s.Require().False(buf.Full())
}

func (s *BoundedTestSuite) TestWriteAndReadTwo() {
	buf := NewBounded(3)
	data := make([]byte, 3)

	n, err := buf.Write([]byte{42, 43, 44})
	s.Require().Equal(3, n)
	s.Require().NoError(err)
	s.Require().False(buf.Empty())
	s.Require().True(buf.Full())

	n, err = buf.Read(data[:2])
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{42, 43}, data[:n])
	s.Require().False(buf.Full())

	n, err = buf.Write([]byte{45, 46})
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().False(buf.Empty())
	s.Require().True(buf.Full())

	n, err = buf.Read(data[:2])
	s.Require().Equal(2, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{44, 45}, data[:n])
	s.Require().False(buf.Empty())
	s.Require().False(buf.Full())

	n, err = buf.Read(data)
	s.Require().Equal(1, n)
	s.Require().NoError(err)
	s.Require().Equal([]byte{46}, data[:n])
	s.Require().True(buf.Empty())
	s.Require().False(buf.Full())
}

func (s *BoundedTestSuite) TestInterface() {
	s.Require().Implements((*Interface)(nil), NewBounded(1))
}

func TestBounded(t *testing.T) {
	suite.Run(t, new(BoundedTestSuite))
}
