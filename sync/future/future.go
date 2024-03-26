package future

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Future holds a result of computation completed in the future.
type Future[V any] struct {
	mu    sync.Mutex
	done  chan struct{}
	error error
	value V
}

// New creates a new incomplete future.
func New[V any]() *Future[V] {
	return &Future[V]{done: make(chan struct{})}
}

// Done returns a channel that's closed when the future is completed.
func (future *Future[V]) Done() <-chan struct{} {
	return future.done
}

// Complete sets a future value and marks it as completed.
// This method should only be called once, subsequent calls will cause a panic.
func (future *Future[V]) Complete(value V) {
	future.mu.Lock()
	defer future.mu.Unlock()

	select {
	case <-future.done:
		panic("future result is already set")
	default:
		future.value = value
		close(future.done)
	}
}

// CompleteUntyped is similar to [Future.Complete] but accepts a value
// of any type as a parameter and panics if it fails to convert to V.
// This method is intended to implement [UntypedInterface].
func (future *Future[V]) CompleteUntyped(untyped any) {
	if value, ok := untyped.(V); ok {
		future.Complete(value)
	} else {
		panic(fmt.Errorf("unexpected future result: expected %T, got %T", value, untyped))
	}
}

// Fail completes a future with an error.
// This method should only be called once, subsequent calls will cause a panic.
func (future *Future[V]) Fail(err error) {
	future.mu.Lock()
	defer future.mu.Unlock()

	select {
	case <-future.done:
		panic("future result is already set")
	default:
		if err == nil {
			panic("future error cannot be nil")
		}
		future.error = err
		close(future.done)
	}
}

// Get returns a future result and a completion status.
func (future *Future[V]) Get() (value V, err error, completed bool) {
	select {
	case <-future.done:
		return future.value, future.error, true
	default:
		return future.makeEmpty(), nil, false
	}
}

// GetUntyped is similar [Future.Get] but returns a value of type any.
// This method is intended to implement [UntypedInterface].
func (future *Future[V]) GetUntyped() (value any, err error, completed bool) {
	return future.Get()
}

// UnmarshalJSON implements a [json.Unmarshaler] interface.
func (future *Future[V]) UnmarshalJSON(data []byte) error {
	value := future.makeEmpty()
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	future.Complete(value)
	return nil
}

func (future *Future[V]) makeEmpty() V {
	var empty V
	return empty
}
