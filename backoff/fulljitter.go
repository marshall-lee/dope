// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import (
	"math/rand"
	"time"
)

// Full jitter algorithm randomizes the delay choosing the value
// between 0 and the exponentially growing window size limited by
// the cap value.
// By randomization it avoids the risk of synchronized retries.
// It may also trigger retries too quickly, which may or may not
// be desirable depending on the case.
//
// Please note that if you choose the base lower than the cap, then
// the algorithm simply won't have a room to grow. So the base will be
// ignored and you simply get a degraded version of the algorithm
// that emits random values in a fixed range.
func FullJitter(base, cap time.Duration) Algorithm {
	current := base
	if current > cap {
		current = cap
	}
	return &fullJitter{current: current, cap: cap}
}

type fullJitter struct {
	current time.Duration
	cap     time.Duration
}

func (alg *fullJitter) Next() time.Duration {
	current := alg.current
	if current < alg.cap {
		if next := current * 2; next < alg.cap {
			alg.current = next
		} else {
			alg.current = alg.cap
		}
	}
	return time.Duration(rand.Int63n(int64(current) + 1))
}
