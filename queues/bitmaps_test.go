// Copyright 2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

var bitmaps []uint32

func init() {
	bitmaps = append(bitmaps, 0)
	for i := 1; i <= 32; i++ {
		x := uint32(1)<<i - 1
		bitmaps = append(bitmaps, x) // 0b1, 0b11, 0b111, 0b1111 ...
		for j := 0; j < 32-i; j++ {
			x <<= 1
			bitmaps = append(bitmaps, x) // 0b00..1..10, 0b1..100, 0b1..1000, ...
		}
	}

	// Some random strings
	bitmaps = append(bitmaps, 1769084016, 1547364371, 828573841)
}
