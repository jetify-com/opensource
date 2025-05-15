package optional

import (
	"encoding/json"
	"fmt"
)

// Option represents an optional value of type T
type Option[T any] struct {
	value T
	set   bool
}

// Some creates a value that is set
func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		set:   true,
	}
}

// None creates an empty value that is not set
func None[T any]() Option[T] {
	return Option[T]{
		set: false,
	}
}

// Ptr creates an Option value from a pointer.
// The value is set if the pointer is not nil, and not set if the pointer is nil.
//
// T should not be a pointer type.
func Ptr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

// Get returns the underlying value and whether it's set
func (o Option[T]) Get() (T, bool) {
	return o.value, o.set
}

// MustGet returns the value or panics if not set
func (o Option[T]) MustGet() T {
	if !o.set {
		panic("option value not set")
	}
	return o.value
}

// GetOrElse returns the contained value or a default if it is not set
func (o Option[T]) GetOrElse(defaultValue T) T {
	if !o.set {
		return defaultValue
	}
	return o.value
}

// IsNone returns true if there is no set value
func (o Option[T]) IsNone() bool {
	return !o.set
}

// String returns a textual representation for debugging/logging
func (o Option[T]) String() string {
	if !o.set {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", o.value)
}

// GoString returns a Go-syntax representation
func (o Option[T]) GoString() string {
	if !o.set {
		return fmt.Sprintf("None[%T]()", o.value)
	}
	return fmt.Sprintf("Some[%T](%#v)", o.value, o.value)
}

// MarshalJSON implements json.Marshaler interface.
// Some(value) marshals to the value's JSON representation
// None() marshals to JSON null
//
// When used with struct fields:
//
//   - None() appears as null by default
//
//   - With 'omitzero' tag: None() values are omitted from JSON
//
//     type Person struct {
//     Name    string       `json:"name"`
//     Age     Option[int]  `json:"age,omitzero"`   // Omitted when None
//     Address Option[string] `json:"address"`      // Appears as null when None
//     }
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.set {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements json.Unmarshaler interface.
// JSON null becomes None(), other values become Some(value)
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = None[T]()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*o = Some(value)
	return nil
}

// IsZero returns 'false' if the value is set, and 'true' if it is not set.
// It makes Option compatible with the 'omitzero' JSON tag.
func (o Option[T]) IsZero() bool {
	return !o.set
}
