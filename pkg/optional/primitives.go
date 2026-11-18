package optional

// Helpers to construct optional values from primitive types.
// A user can always use Some(value) to construct an optional value
// even for primitive types:
//
//	opt := optional.Some(42)
//	opt := optional.Some[float64](42)
//
// But the constructors here sometimes make the intent more clear even
// if they are equivalent:
//
//	opt := optional.Int(42)
//	opt := optional.Float64(42)

// Int returns an Option wrapping the given int value.
func Int(v int) Option[int] {
	return Some(v)
}

// Int8 returns an Option wrapping the given int8 value.
func Int8(v int8) Option[int8] {
	return Some(v)
}

// Int16 returns an Option wrapping the given int16 value.
func Int16(v int16) Option[int16] {
	return Some(v)
}

// Int32 returns an Option wrapping the given int32 value.
func Int32(v int32) Option[int32] {
	return Some(v)
}

// Int64 returns an Option wrapping the given int64 value.
func Int64(v int64) Option[int64] {
	return Some(v)
}

// Uint returns an Option wrapping the given uint value.
func Uint(v uint) Option[uint] {
	return Some(v)
}

// Uint8 returns an Option wrapping the given uint8 value.
func Uint8(v uint8) Option[uint8] {
	return Some(v)
}

// Uint16 returns an Option wrapping the given uint16 value.
func Uint16(v uint16) Option[uint16] {
	return Some(v)
}

// Uint32 returns an Option wrapping the given uint32 value.
func Uint32(v uint32) Option[uint32] {
	return Some(v)
}

// Uint64 returns an Option wrapping the given uint64 value.
func Uint64(v uint64) Option[uint64] {
	return Some(v)
}

// Float32 returns an Option wrapping the given float32 value.
func Float32(v float32) Option[float32] {
	return Some(v)
}

// Float64 returns an Option wrapping the given float64 value.
func Float64(v float64) Option[float64] {
	return Some(v)
}

// Complex64 returns an Option wrapping the given complex64 value.
func Complex64(v complex64) Option[complex64] {
	return Some(v)
}

// Complex128 returns an Option wrapping the given complex128 value.
func Complex128(v complex128) Option[complex128] {
	return Some(v)
}

// Bool returns an Option wrapping the given bool value.
func Bool(v bool) Option[bool] {
	return Some(v)
}

// String returns an Option wrapping the given string value.
func String(v string) Option[string] {
	return Some(v)
}

// Byte returns an Option wrapping the given byte value.
func Byte(v byte) Option[byte] {
	return Some(v)
}

// Rune returns an Option wrapping the given rune value.
func Rune(v rune) Option[rune] {
	return Some(v)
}
