// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDecorr(t *testing.T) {
	base := 100 * time.Millisecond
	cap := 500 * time.Millisecond
	backoff := NewDecorr(base, cap)
	for i := 0; i < 100; i++ {
		next := backoff.Next()
		require.GreaterOrEqual(t, next, base)
		require.LessOrEqual(t, next, cap)
	}
}

func TestDecorrBiggerBase(t *testing.T) {
	base := 500 * time.Millisecond
	cap := 100 * time.Millisecond
	backoff := NewDecorr(base, cap)
	for i := 0; i < 100; i++ {
		require.Equal(t, cap, backoff.Next())
	}
}
