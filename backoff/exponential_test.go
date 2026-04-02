// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExponential(t *testing.T) {
	backoff := Exponential(100*time.Millisecond, time.Second)
	require.Equal(t, 100*time.Millisecond, backoff.Next())
	require.Equal(t, 200*time.Millisecond, backoff.Next())
	require.Equal(t, 400*time.Millisecond, backoff.Next())
	require.Equal(t, 800*time.Millisecond, backoff.Next())
	require.Equal(t, 1000*time.Millisecond, backoff.Next())
	require.Equal(t, 1000*time.Millisecond, backoff.Next())
}

func TestExponentialBiggerBase(t *testing.T) {
	base := 200 * time.Millisecond
	cap := 100 * time.Millisecond
	backoff := Exponential(base, cap)
	for i := 0; i < 10; i++ {
		require.Equal(t, cap, backoff.Next())
	}
}
