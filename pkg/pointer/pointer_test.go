package pointer

import "testing"

func check[T comparable](t *testing.T, name string, p *T, want T) {
	t.Helper()
	if p == nil {
		t.Fatalf("%s() returned nil", name)
	}
	if *p != want {
		t.Errorf("*%s() = %v, want %v", name, *p, want)
	}
}

func TestConstructors(t *testing.T) {
	check(t, "Ptr", Ptr(struct{ A int }{1}), struct{ A int }{1})
	check(t, "Int", Int(-1), -1)
	check(t, "Int8", Int8(-8), int8(-8))
	check(t, "Int16", Int16(-16), int16(-16))
	check(t, "Int32", Int32(-32), int32(-32))
	check(t, "Int64", Int64(-64), int64(-64))
	check(t, "Uint", Uint(1), uint(1))
	check(t, "Uint8", Uint8(8), uint8(8))
	check(t, "Uint16", Uint16(16), uint16(16))
	check(t, "Uint32", Uint32(32), uint32(32))
	check(t, "Uint64", Uint64(64), uint64(64))
	check(t, "Float32", Float32(1.5), float32(1.5))
	check(t, "Float64", Float64(2.5), 2.5)
	check(t, "Complex64", Complex64(1+2i), complex64(1+2i))
	check(t, "Complex128", Complex128(3+4i), 3+4i)
	check(t, "Bool", Bool(true), true)
	check(t, "String", String("s"), "s")
	check(t, "Byte", Byte('b'), byte('b'))
	check(t, "Rune", Rune('r'), 'r')
}

func TestConstructorsReturnFreshPointers(t *testing.T) {
	v := 5
	p := Ptr(v)
	if p == &v {
		t.Fatal("Ptr() returned the address of its argument's caller variable")
	}
	*p = 6
	if v != 5 {
		t.Errorf("mutating *Ptr(v) changed v to %d", v)
	}
	if a, b := Int(1), Int(1); a == b {
		t.Error("Int() returned the same pointer twice")
	}
}

func TestValueOr(t *testing.T) {
	if got := ValueOr[int](nil, 7); got != 7 {
		t.Errorf("ValueOr(nil, 7) = %d, want 7", got)
	}
	if got := ValueOr(Int(3), 7); got != 3 {
		t.Errorf("ValueOr(&3, 7) = %d, want 3", got)
	}
	// A pointer to the zero value is not nil, so the default is not used.
	if got := ValueOr(Int(0), 7); got != 0 {
		t.Errorf("ValueOr(&0, 7) = %d, want 0", got)
	}
}

func TestValue(t *testing.T) {
	if got := Value[string](nil); got != "" {
		t.Errorf("Value(nil) = %q, want empty", got)
	}
	if got := Value(String("x")); got != "x" {
		t.Errorf("Value(&x) = %q, want x", got)
	}
	type s struct{ N int }
	if got := Value[s](nil); got != (s{}) {
		t.Errorf("Value(nil) = %+v, want zero struct", got)
	}
}
