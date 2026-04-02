// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backoff

import (
	"time"
)

type constant struct {
	base time.Duration
}

func (alg constant) Next() time.Duration {
	return alg.base
}
