// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzBoundedPushPop(f *testing.F) {
	for _, capacity := range []int{1, 2, 4, 8, 16, 32} {
		for _, x := range bitmaps {
			f.Add(capacity, x)
		}
	}
	f.Fuzz(func(t *testing.T, capacity int, x uint32) {
		queue := NewBounded[int](capacity)
		var invariant []int
		for i := 1; i <= 32; i++ {
			if x&1 == 1 {
				ok := queue.Push(i)
				if len(invariant) < capacity {
					require.True(t, ok)
					invariant = append(invariant, i)
				} else {
					require.False(t, ok)
				}

			} else {
				value, ok := queue.Pop()
				if len(invariant) != 0 {
					require.True(t, ok)
					require.Equal(t, invariant[0], value)
					invariant = invariant[1:]
				} else {
					require.False(t, ok)
				}
			}

			require.Equal(t, len(invariant) == 0, queue.Empty())
			require.Equal(t, len(invariant) == capacity, queue.Full())
			require.Equal(t, len(invariant), queue.Len())

			x >>= 1
		}
	})
}

func FuzzBoundedPushSomePopSome(f *testing.F) {
	for _, capacity := range []int{1, 2, 4, 8, 16, 32} {
		for _, x := range bitmaps {
			f.Add(capacity, x)
		}
	}
	f.Fuzz(func(t *testing.T, capacity int, x uint32) {
		queue := NewBounded[int](capacity)
		invariant := make([]int, 0)
		for x > 0 {
			var data []int
			for i := 1; x&1 == 1; i++ {
				data = append(data, i)
				x >>= 1
			}
			n := queue.PushSome(data)
			invariant = append(invariant, data...)
			if overflow := len(invariant) - capacity; overflow > 0 {
				invariant = invariant[:capacity]
				require.Equal(t, len(data)-overflow, n)
			} else {
				require.Equal(t, len(data), n)
			}

			if x > 0 {
				nPop := 0
				for x&1 == 0 {
					nPop += 1
					x >>= 1
				}
				out := make([]int, nPop)
				n = queue.PopSome(out)
				if nPop > len(invariant) {
					nPop = len(invariant)
				}
				require.Equal(t, invariant[:nPop], out[:n])
				invariant = invariant[nPop:]
			}

			require.Equal(t, len(invariant) == 0, queue.Empty())
			require.Equal(t, len(invariant) == capacity, queue.Full())
			require.Equal(t, len(invariant), queue.Len())
		}
	})
}
