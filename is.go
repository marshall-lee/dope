// Copyright 2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dope

// Is checks that a value is of a given (generic) type.
func Is[T any](val any) bool {
	_, ok := val.(T)
	return ok
}

// IsEmpty checks that a value is equal to an empty value of a given comparable (generic) type.
func IsEmpty[T comparable](val T) bool {
	var empty T
	return val == empty
}
