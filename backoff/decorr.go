// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import (
	"math/rand"
	"time"
)

// Decorrelated jitter algorithm generates the next value
// depending on a previous step. It generates a random value
// between base and 3 * prev but no more than the cap value.
// This algorithm is the best in avoiding the risk of synchronous
// retries and it also avoids too short retries.
//
// Please note that in practice you should choose the base lower than
// the cap. Otherwise, the algorithm would emit constant values which is
// probably not what you wanted.
func NewDecorr(base, cap time.Duration) Algorithm {
	if base >= cap {
		return constant{base: cap}
	}
	return &decorr{base: base, current: base, cap: cap}
}

type decorr struct {
	base    time.Duration
	current time.Duration
	cap     time.Duration
}

func (alg *decorr) Next() time.Duration {
	current := time.Duration(int64(alg.base) + rand.Int63n(int64(alg.current)*3-int64(alg.base)+1))
	if current > alg.cap {
		current = alg.cap
	}
	alg.current = current
	return time.Duration(current)
}
