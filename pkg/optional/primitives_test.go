package optional

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertSome[T comparable](t *testing.T, o Option[T], want T) {
	t.Helper()
	got, ok := o.Get()
	assert.True(t, ok)
	assert.Equal(t, want, got)
	assert.False(t, o.IsNone())
}

func TestPrimitiveConstructors(t *testing.T) {
	assertSome(t, Int(-1), -1)
	assertSome(t, Int8(-8), int8(-8))
	assertSome(t, Int16(-16), int16(-16))
	assertSome(t, Int32(-32), int32(-32))
	assertSome(t, Int64(-64), int64(-64))
	assertSome(t, Uint(1), uint(1))
	assertSome(t, Uint8(8), uint8(8))
	assertSome(t, Uint16(16), uint16(16))
	assertSome(t, Uint32(32), uint32(32))
	assertSome(t, Uint64(64), uint64(64))
	assertSome(t, Float32(1.5), float32(1.5))
	assertSome(t, Float64(2.5), 2.5)
	assertSome(t, Complex64(1+2i), complex64(1+2i))
	assertSome(t, Complex128(3+4i), 3+4i)
	assertSome(t, Bool(false), false)
	assertSome(t, String(""), "")
	assertSome(t, Byte('b'), byte('b'))
	assertSome(t, Rune('r'), 'r')
}

func TestPrimitiveZeroValuesAreSome(t *testing.T) {
	// Wrapping a zero value must still produce a present option.
	assert.False(t, Int(0).IsNone())
	assert.False(t, Bool(false).IsNone())
	assert.False(t, String("").IsNone())
	assert.False(t, Int(0).IsZero())
}
