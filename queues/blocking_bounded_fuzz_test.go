// Copyright 2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"slices"
	"testing"
	"time"

	internal "github.com/marshall-lee/dope/internal/tests"
	"github.com/stretchr/testify/require"
)

func FuzzBlockingBoundedPushPop(f *testing.F) {
	for _, capacity := range []int{1, 2, 4, 8, 16, 32} {
		for _, x := range bitmaps {
			f.Add(capacity, x)
		}
	}
	f.Fuzz(func(t *testing.T, capacity int, x uint32) {
		queue := NewBlockingBounded[int](capacity)
		invariant := make(chan int, capacity)

		type popResult struct {
			value int
			ok    bool
		}

		var (
			pushs, pushsInv []internal.GoroutineCapture[struct{}]
			pops, popsInv   []internal.GoroutineCapture[popResult]
			depth           int
		)

		for i := 1; i <= 32; i++ {
			if x&1 == 1 {
				push := internal.GoCapture(func() { queue.Push(i) })
				pushInv := internal.GoCapture(func() { invariant <- i })

				if isFull := depth >= capacity; !isFull {
					require.Eventually(t, push.IsDone, 50*time.Millisecond, 10*time.Microsecond, "Push did not return")
					require.Eventually(t, pushInv.IsDone, 50*time.Millisecond, 10*time.Microsecond, "chan<- did not return")
				}

				pushs = append(pushs, push)
				pushsInv = append(pushsInv, pushInv)
				depth += 1
			} else {
				pop := internal.GoCaptureWithReturnValue(func() popResult { value, ok := queue.Pop(); return popResult{value, ok} })
				popInv := internal.GoCaptureWithReturnValue(func() popResult { value, ok := <-invariant; return popResult{value, ok} })

				if isEmpty := depth <= 0; !isEmpty {
					require.Eventually(t, pop.IsDone, 50*time.Millisecond, 10*time.Microsecond, "Pop did not return")
					require.Eventually(t, popInv.IsDone, 50*time.Millisecond, 10*time.Microsecond, "<-chan did not return")
					require.True(t, pop.Val().ok)
					require.True(t, popInv.Val().ok)
				}

				pops = append(pops, pop)
				popsInv = append(popsInv, popInv)
				depth -= 1
			}
			x >>= 1
		}
		var (
			popVals       []int
			stuckPops     []internal.GoroutineCapture[popResult]
			stuckPopsInv  []internal.GoroutineCapture[popResult]
			stuckPushs    []internal.GoroutineCapture[struct{}]
			stuckPushsInv []internal.GoroutineCapture[struct{}]
		)

		for _, push := range pushs {
			if !push.IsDone() {
				stuckPushs = append(stuckPushs, push)
			}

		}
		for _, push := range pushsInv {
			if !push.IsDone() {
				stuckPushsInv = append(stuckPushsInv, push)
			}
		}

		for _, pop := range pops {
			if !pop.IsDone() {
				stuckPops = append(stuckPops, pop)
				continue
			}
			v := pop.Val()
			require.True(t, v.ok)
			require.GreaterOrEqual(t, v.value, 1)
			require.LessOrEqual(t, v.value, 32)

			popVals = append(popVals, v.value)
		}
		for _, pop := range popsInv {
			if !pop.IsDone() {
				stuckPopsInv = append(stuckPopsInv, pop)
				continue
			}
			v := pop.Val()
			require.True(t, v.ok)
			require.GreaterOrEqual(t, v.value, 1)
			require.LessOrEqual(t, v.value, 32)
		}

		require.Equal(t, len(stuckPopsInv), len(stuckPops))
		if depth < 0 {
			require.Equal(t, -depth, len(stuckPops))
		} else {
			require.Equal(t, 0, len(stuckPops))
		}

		require.Equal(t, len(stuckPushsInv), len(stuckPushs))
		if depth > capacity {
			require.Equal(t, depth-capacity, len(stuckPushs))
		} else {
			require.Equal(t, 0, len(stuckPushs))
		}

		slices.Sort(popVals)
		for i, val := range popVals {
			if i > 0 {
				require.Greater(t, val, popVals[i-1]) // Ensure no duplicates.
			}
		}

		queue.Close()
		close(invariant)

		for _, pop := range stuckPops {
			require.Eventually(t, pop.IsDone, 50*time.Millisecond, 10*time.Microsecond)
			require.False(t, pop.Val().ok)
		}
		for _, pop := range stuckPopsInv {
			require.Eventually(t, pop.IsDone, 50*time.Millisecond, 10*time.Microsecond)
			require.False(t, pop.Val().ok)
		}
		for _, push := range stuckPushs {
			require.Eventually(t, push.IsPanicked, 50*time.Millisecond, 10*time.Microsecond)
		}
		for _, push := range stuckPushsInv {
			require.Eventually(t, push.IsPanicked, 50*time.Millisecond, 10*time.Microsecond)
		}
	})
}
