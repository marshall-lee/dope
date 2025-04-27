// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package futures

// UntypedInterface is useful when for example you need a collection of [Future]
// objects with different type parameters.
type UntypedInterface interface {
	// Done returns a channel that's closed when the future is completed.
	Done() <-chan struct{}
	// Complete sets a future value and marks it as completed.
	CompleteUntyped(value any)
	// Fail completes a future with an error.
	Fail(err error)
	// Get returns a future result and a completion status.
	GetUntyped() (value any, err error, ok bool)
}

// Untyped future can be completed with a value of any type.
type Untyped Future[any]

// NewUntyped creates a new incomplete [Untyped] future.
func NewUntyped() *Untyped {
	return (*Untyped)(New[any]())
}

// Done returns a channel that's closed when the future is completed.
func (untyped *Untyped) Done() <-chan struct{} {
	return untyped.future().Done()
}

// CompleteUntyped sets a future value and marks it as completed.
// This method should only be called once, subsequent calls will cause a panic.
func (untyped *Untyped) CompleteUntyped(value any) {
	untyped.future().Complete(value)
}

// Fail completes a future with an error.
// This method should only be called once, subsequent calls will cause a panic.
func (untyped *Untyped) Fail(err error) {
	untyped.future().Fail(err)
}

// GetUntyped returns a future result and a completion status.
func (untyped *Untyped) GetUntyped() (value any, err error, ok bool) {
	return untyped.future().Get()
}

func (untyped *Untyped) future() *Future[any] {
	return (*Future[any])(untyped)
}
