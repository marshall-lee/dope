// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"errors"
)

var (
	ErrCloseClosed = errors.New("queues: close of a closed queue")
	ErrPushClosed = errors.New("queues: push to a closed queue")
)
