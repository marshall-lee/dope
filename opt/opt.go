// Copyright 2026 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

// Setter is an interface of at option setter. Each setter is intended
// to modify in-place a specific part of an object given by its pointer.
type Setter[T any] interface {
	ApplyOption(*T)
}

// Apply sequentially invokes a list of setters on a given object.
func Apply[T any](obj *T, setters ...Setter[T]) {
	for _, setter := range setters {
		setter.ApplyOption(obj)
	}
}

type applyFunc[T any] func(*T)

func (f applyFunc[T]) ApplyOption(obj *T) {
	f(obj)
}

// ApplyFunc wraps a Go function into a Setter.
func ApplyFunc[T any](fn func(*T)) Setter[T] {
	return applyFunc[T](fn)
}
