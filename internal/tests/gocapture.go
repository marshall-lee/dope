// Copyright 2024-2025 Vladimir Kochnev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queues

import (
	"runtime"

	"github.com/marshall-lee/dope/sync/futures"
)

type GoroutineCapture[T any] struct {
	future *futures.Future[goroutineCaptureInner[T]]
}

type goroutineCaptureInner[T any] struct {
	panicked bool
	err      any
	val      T
}

func GoCapture(f func()) GoroutineCapture[struct{}] {
	return GoCaptureWithReturnValue(func() (void struct{}) { f(); return void })
}

func GoCaptureWithReturnValue[T any](f func() T) GoroutineCapture[T] {
	future := futures.New[goroutineCaptureInner[T]]()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				future.Complete(goroutineCaptureInner[T]{panicked: true, err: err})
			}
		}()

		val := f()

		future.Complete(goroutineCaptureInner[T]{val: val})
	}()
	runtime.Gosched()

	return GoroutineCapture[T]{future}
}

func (cap GoroutineCapture[T]) Done() <-chan struct{} {
	return cap.future.Done()
}

func (cap GoroutineCapture[T]) IsDone() bool {
	select {
	case <-cap.future.Done():
		return true
	default:
		return false
	}
}

func (cap GoroutineCapture[T]) IsPanicked() bool {
	res, err, _ := cap.future.Get()
	if err != nil {
		panic(err)
	}
	return res.panicked
}

func (cap GoroutineCapture[T]) Val() T {
	res, err, _ := cap.future.Get()
	if err != nil {
		panic(err)
	}
	return res.val
}

func (cap GoroutineCapture[T]) Err() any {
	res, err, _ := cap.future.Get()
	if err != nil {
		panic(err)
	}
	return res.err
}
