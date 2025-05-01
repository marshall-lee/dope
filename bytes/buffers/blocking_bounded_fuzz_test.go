// Copyright 2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buffers

import (
	"errors"
	"io"
	"testing"
	"time"

	internal "github.com/marshall-lee/dope/internal/tests"
	"github.com/stretchr/testify/require"
)

func FuzzBlockingBoundedSingleWriteReadLoop(f *testing.F) {
	strs := []string{
		"",
		"a",
		"abc",
		"foobarbaz",
		"abcdefghijklmnopqrstuvwxyz",
		"0123456789abcdefghijklmnopqrstuvwxyz",
		"абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
	}
	caps := []int{1, 2, 4, 8, 16, 32}
	for _, readCapacity := range caps {
		for _, bufCapacity := range caps {
			for _, str := range strs {
				f.Add(readCapacity, bufCapacity, str)
			}
		}
	}
	f.Fuzz(func(t *testing.T, readCapacity, bufCapacity int, str string) {
		buf := NewBlockingBounded(bufCapacity)

		var invariant string
		consumer := internal.GoCapture(func() {
			data := make([]byte, readCapacity)
			for {
				n, err := buf.Read(data)
				if n > 0 {
					invariant = invariant + string(data[:n])
				}

				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					panic(err)
				}
			}
		})

		writeData := []byte(str)
		n, err := buf.Write(writeData)
		require.NoError(t, err)
		require.Equal(t, len(writeData), n)

		err = buf.Close()
		require.NoError(t, err)

		require.Eventually(t, consumer.IsDone, 50*time.Millisecond, 10*time.Microsecond)
		require.False(t, consumer.IsPanicked())
		require.Equal(t, str, invariant)
	})
}

func FuzzBlockingBoundedPartWriteReadLoop(f *testing.F) {
	strs := []string{
		"",
		"a",
		"abc",
		"foobarbaz",
		"abcdefghijklmnopqrstuvwxyz",
		"0123456789abcdefghijklmnopqrstuvwxyz",
		"абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
	}
	caps := []int{1, 2, 4, 8, 16, 32}
	for _, writeSize := range caps {
		for _, readCapacity := range caps {
			for _, bufCapacity := range caps {
				for _, str := range strs {
					f.Add(writeSize, readCapacity, bufCapacity, str)
				}
			}
		}
	}
	f.Fuzz(func(t *testing.T, writeSize, readCapacity, bufCapacity int, str string) {
		buf := NewBlockingBounded(bufCapacity)

		var invariant string
		consumer := internal.GoCapture(func() {
			data := make([]byte, readCapacity)
			for {
				n, err := buf.Read(data)
				if n > 0 {
					invariant = invariant + string(data[:n])
				}

				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					panic(err)
				}
			}
		})

		writeData := []byte(str)

		for i := 0; i < len(writeData); i += writeSize {
			m := writeSize
			if i+m > len(writeData) {
				m = len(writeData) - i
			}
			n, err := buf.Write(writeData[i : i+m])
			require.NoError(t, err)
			require.Equal(t, m, n)
		}

		err := buf.Close()
		require.NoError(t, err)

		require.Eventually(t, consumer.IsDone, 50*time.Millisecond, 10*time.Microsecond)
		require.False(t, consumer.IsPanicked())
		require.Equal(t, str, invariant)
	})
}
