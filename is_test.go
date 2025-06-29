// Copyright 2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dope

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIs(t *testing.T) {
	type mystruct struct{}
	type anotherstruct struct{}

	require.True(t, Is[string]("str"))
	require.True(t, Is[int](123))

	require.False(t, Is[int]("str"))
	require.False(t, Is[string](123))

	require.True(t, Is[mystruct](mystruct{}))
	require.False(t, Is[mystruct](anotherstruct{}))
}
