// Package optional provides a generic Option[T] type and helpers for representing
// values that may or may not be present.
//
// The Option[T] type is a lightweight, allocation-free abstraction similar to the
// "Maybe" or "Optional" types found in other languages.  It is useful for
// expressing that a value might be absent without resorting to pointers or
// sentinel values.
//
// # Basic construction
//
// There are two canonical constructors:
//
//	opt := optional.Some(42) // value is set
//	opt := optional.None[int]() // value is not set
//
// Accessing values
//
//	v, ok := opt.Get()         // returns value and whether it is set
//	v := opt.GetOrElse(0)      // returns value or a default
//	v := opt.MustGet()         // panics if value is not set
//	if opt.IsNone() { ... }    // tests for absence
//
// Pointers
//
//	opt := optional.Ptr(ptr)   // wraps a *T, set to the value of ptr when ptr != nil
//
// # Primitive helpers
//
// For convenience the package exposes helper constructors for primitive types:
//
//	optional.Int(1)        // Option[int]
//	optional.Float64(3.14) // Option[float64]
//	optional.String("hi") // Option[string]
//
// These are entirely equivalent to calling Some(value) but sometimes improve
// readability.
//
// # JSON integration
//
// Option implements json.Marshaler and json.Unmarshaler.  A None value
// marshals to JSON null, while Some(value) marshals to the JSON representation
// of the contained value.
//
// When Option fields are used inside structs the 'omitzero' tag can be used to
// omit None() values entirely:
//
//	type Person struct {
//	    Name string          `json:"name"`
//	    Age  optional.Option[int] `json:"age,omitzero"` // omitted when None
//	}
//
// # Zero value behavior
//
// The zero value of Option[T] is considered "not set" (None).  As a result an
// Option field in a struct that is left uninitialized will behave as None().
//
// # Concurrency
//
// Option is a value type. Copies are immutable and safe for concurrent reads.
// The only mutator is (*Option).UnmarshalJSON; if you share the same Option
// *pointer* between goroutines, guard that pointer with a mutex while
// unmarshalling.
package optional
