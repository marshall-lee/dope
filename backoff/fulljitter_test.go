// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFullJitter(t *testing.T) {
	base := time.Millisecond
	backoff := FullJitter(base, time.Second)
	exp := base
	for i := 0; i < 10; i++ {
		next := backoff.Next()
		require.GreaterOrEqual(t, next, time.Duration(0))
		require.LessOrEqual(t, next, exp)
		exp *= 2
	}
}

func TestFullJitterBiggerBase(t *testing.T) {
	base := 200 * time.Millisecond
	cap := 100 * time.Millisecond
	backoff := FullJitter(base, cap)
	for i := 0; i < 10; i++ {
		next := backoff.Next()
		require.GreaterOrEqual(t, next, time.Duration(0))
		require.LessOrEqual(t, next, cap)
	}
}
