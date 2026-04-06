// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import "time"

// Exponential backoff is a well-known algorithm.
// Starting with base, each next delay is twice long but
// is no more than the cap value.
//
// Please note that in practice you should choose the base lower than
// the cap. Otherwise, the algorithm would emit constant values which is
// probably not what you wanted.
func NewExponential(base, cap time.Duration) Algorithm {
	if base >= cap {
		return constant{base: cap}
	}
	return &exponential{current: base, cap: cap}
}

type exponential struct {
	current time.Duration
	cap     time.Duration
}

func (alg *exponential) Next() time.Duration {
	current := alg.current
	if current < alg.cap {
		if next := current * 2; next < alg.cap {
			alg.current = next
		} else {
			alg.current = alg.cap
		}
	}
	return current
}
