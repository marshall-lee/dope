// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"sync"
)

func GoCapturePanic(f func(), locker sync.Locker, panickedPtr *bool, errPtr *any) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				locker.Lock()
				*panickedPtr = true
				*errPtr = err
				locker.Unlock()
			}
		}()
		f()
	}()
}
