package futures

import (
	"testing"
	"encoding/json"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCompleteAndGet(t *testing.T) {
	future := New[int]()
	future.Complete(42)
	value, err, ok := future.Get()
	require.Equal(t, 42, value)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestGetIncomplete(t *testing.T) {
	future := New[int]()
	_, err, ok := future.Get()
	require.NoError(t, err)
	require.False(t, ok)
}

func TestDone(t *testing.T) {
	future := New[int]()
	select {
	case <-future.Done():
		require.Fail(t, "future must not be done yet")
	default:
	}

	const sleepDuration = time.Millisecond * 100
	start := time.Now()
	go func() {
		time.Sleep(sleepDuration)
		future.Complete(42)
	}()
	<-future.Done()
	require.InEpsilon(t, sleepDuration, time.Since(start), 0.05)
	value, err, ok := future.Get()
	require.Equal(t, 42, value)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestUnmarshalJSON(t *testing.T) {
	future := New[struct{
		Foo int `json:"foo"`
		Bar string `json:"bar"`
	}]()

	err := json.Unmarshal([]byte(`{"foo":42,"bar":"baz"}`), future)
	require.NoError(t, err)

	value, err, ok := future.Get()
	require.Equal(t, 42, value.Foo)
	require.Equal(t, "baz", value.Bar)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestUntypedInterface(t *testing.T) {
	require.Implements(t, (*UntypedInterface)(nil), New[int]())
	require.Implements(t, (*UntypedInterface)(nil), New[string]())
	require.Implements(t, (*UntypedInterface)(nil), New[struct{}]())
	require.Implements(t, (*UntypedInterface)(nil), New[any]())
	require.Implements(t, (*UntypedInterface)(nil), NewUntyped())
}
