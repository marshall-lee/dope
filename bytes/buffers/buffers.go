// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buffers

import "io"

type Interface interface {
	io.Writer
	io.ByteWriter
	io.Reader
	io.ByteReader
}
