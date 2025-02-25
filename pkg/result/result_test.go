package result

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstructors(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		r := Ok(42)
		assert.True(t, r.IsOk())
		assert.Equal(t, 42, r.MustGet())
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New("test error")
		r := Err[int](err)
		assert.True(t, r.IsErr())
		// We don't want to use the errorlint rule here because we are actually
		// testing that Err returns the exact error we set.
		//nolint:errorlint
		assert.Equal(t, err, r.Err())
	})

	t.Run("Errf", func(t *testing.T) {
		r := Errf[int]("error %d", 42)
		assert.True(t, r.IsErr())
		assert.EqualError(t, r.Err(), "error 42")
	})

	t.Run("From", func(t *testing.T) {
		r1 := From(42, nil)
		assert.True(t, r1.IsOk())

		err := errors.New("test error")
		r2 := From(42, err)
		assert.True(t, r2.IsErr())
	})
}

func TestPredicates(t *testing.T) {
	t.Run("IsOk and IsErr", func(t *testing.T) {
		ok := Ok(42)
		assert.True(t, ok.IsOk())
		assert.False(t, ok.IsErr())

		err := Err[int](errors.New("test"))
		assert.True(t, err.IsErr())
		assert.False(t, err.IsOk())
	})
}

func TestMethods(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		v, err := Ok(42).Get()
		assert.NoError(t, err)
		assert.Equal(t, 42, v)

		testErr := errors.New("test")
		v, err = Err[int](testErr).Get()
		// We don't want to use the errorlint rule here because we want to test
		// that Get returns the exact error we set.
		//nolint:errorlint
		assert.Equal(t, testErr, err)
		assert.Zero(t, v)
	})

	t.Run("MustGet", func(t *testing.T) {
		v := Ok(42).MustGet()
		assert.Equal(t, 42, v)

		assert.Panics(t, func() {
			Err[int](errors.New("test")).MustGet()
		})
	})

	t.Run("OrElse", func(t *testing.T) {
		assert.Equal(t, 42, Ok(42).OrElse(10))
		assert.Equal(t, 10, Err[int](errors.New("test")).OrElse(10))
	})
}

func TestFormatting(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "Ok(42)", Ok(42).String())
		assert.Equal(t, "Err(test)", Err[int](errors.New("test")).String())
	})

	t.Run("GoString", func(t *testing.T) {
		assert.Equal(t, `Ok[int](42)`, Ok(42).GoString())
		assert.Equal(t, `Err[int]("test")`, Err[int](errors.New("test")).GoString())
	})
}

func TestDo(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		r := Do(func() int { return 42 })
		assert.True(t, r.IsOk())
		assert.Equal(t, 42, r.MustGet())
	})

	t.Run("panic with error", func(t *testing.T) {
		r := Do(func() int {
			panic(errors.New("test panic"))
		})
		assert.True(t, r.IsErr())
		assert.EqualError(t, r.Err(), "test panic")
	})

	t.Run("panic with string", func(t *testing.T) {
		r := Do(func() int {
			panic("test panic")
		})
		assert.True(t, r.IsErr())
		assert.EqualError(t, r.Err(), "test panic")
	})
}

func TestGo(t *testing.T) {
	t.Run("successful async execution", func(t *testing.T) {
		ch := Go(func() (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 42, nil
		})

		r := <-ch
		assert.True(t, r.IsOk())
		assert.Equal(t, 42, r.MustGet())
	})

	t.Run("error async execution", func(t *testing.T) {
		ch := Go(func() (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 0, fmt.Errorf("async error")
		})

		r := <-ch
		assert.True(t, r.IsErr())
		assert.EqualError(t, r.Err(), "async error")
	})

	t.Run("channel closes after result", func(t *testing.T) {
		ch := Go(func() (int, error) {
			return 42, nil
		})

		<-ch // Read the result
		_, ok := <-ch
		assert.False(t, ok, "Expected channel to be closed after result")
	})
}
