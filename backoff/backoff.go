// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import (
	"context"
	"time"
)

// Backoff is the algorithm wrapper that implements actual sleeping
// methods.
type Backoff struct {
	alg Algorithm
}

// Algorithm represents a backoff algorithm.
type Algorithm interface {
	Next() time.Duration
}

func New(alg Algorithm) Backoff {
	return Backoff{alg: alg}
}

// Sleep pauses the current goroutine for a duration determined by the algorithm.
func (b Backoff) Sleep() {
	time.Sleep(b.alg.Next())
}

// SleepWithContext pauses the current goroutine but can exit earlier if ctx
// is canceled or deadlined.
func (b Backoff) SleepWithContext(ctx context.Context) error {
	timer := time.NewTimer(b.alg.Next())
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// After waits for a duration determined by the algorithm,
// then sends the current local time on the returned channel.
func (b Backoff) After() <-chan time.Time {
	return time.After(b.alg.Next())
}
