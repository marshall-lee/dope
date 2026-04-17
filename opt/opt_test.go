// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type multiplierT struct{}

func (_ multiplierT) ApplyOption(val *int) {
	*val = *val + *val
}

var (
	adder = ApplyFunc(func(val *int) {
		*val += 1
	})
	multiplier = multiplierT{}
)

func TestApply(t *testing.T) {
	val := 0
	Apply(&val, adder)
	require.Equal(t, 1, val)
	Apply(&val, adder, adder, adder)
	require.Equal(t, 4, val)
	Apply(&val, multiplier, adder, multiplier)
	require.Equal(t, 18, val)
}
