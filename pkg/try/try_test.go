package try

import (
	"errors"
	"testing"

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

	t.Run("GetOrElse", func(t *testing.T) {
		assert.Equal(t, 42, Ok(42).GetOrElse(10))
		assert.Equal(t, 10, Err[int](errors.New("test")).GetOrElse(10))
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
