package pointer

// Package pointer offers a collection of helper functions for working with Go pointers in a
// type-safe and ergonomic way. The utilities fall into three broad categories:
//
//   1. Pointer Constructors
//      For each of Go's built-in numeric, string and boolean types there is a constructor
//      that takes a value and returns a pointer to that value. These helpers eliminate the
//      boilerplate of creating a local variable and taking its address when you need to pass
//      a pointer literal, e.g.
//
//          user.Age = pointer.Int(42)   // instead of `v := 42; user.Age = &v`
//
//      The type-specific constructors (Int, String, Bool, etc.) exist primarily for code clarity
//      and readability. The generic Ptr[T] function can be used with any type, including
//      primitives (e.g., pointer.Ptr(42)), but the specific constructors can make the intent
//      clearer without requiring type parameters:
//
//          // Both of these work, but the first is more readable
//          score := pointer.Float64(42)         // More idiomatic
//          score := pointer.Ptr[float64](42)    // Works, but more verbose
//
//      The constructors are especially useful when dealing with structs that have pointer fields:
//
//          type User struct {
//              ID       *string
//              Name     *string
//              Age      *int
//              Score    *float64
//              Active   *bool
//          }
//
//          // Initialize a struct with pointer fields
//          user := User{
//              ID:       pointer.String("user123"),
//              Name:     pointer.String("Alice"),
//              Age:      pointer.Int(30),
//              Score:    pointer.Float64(98.5),
//              Active:   pointer.Bool(true),
//          }
//
//   2. Safe Dereferencing
//      Value[T] returns the dereferenced value of *T, falling back to the zero value of the
//      element type when the pointer is nil. This makes it convenient to read optional
//      pointer fields without additional nil checks:
//
//          // Safe even if userAge is nil
//          age := pointer.Value(userAge) // returns 0 if nil
//
//          // Equivalent to:
//          // var age int
//          // if userAge != nil {
//          //     age = *userAge
//          // }
//
//   3. Defaulting Helpers
//      ValueOr[T] behaves like Value but lets the caller supply an explicit default value to
//      use when the pointer is nil:
//
//          // Returns the dereferenced value if not nil, or 18 if nil
//          age := pointer.ValueOr(userAge, 18)
//
//          // Equivalent to:
//          // age := 18
//          // if userAge != nil {
//          //     age = *userAge
//          // }
//
// All helpers are implemented as small, inlinable one-liners so they impose no runtime cost.
// They are intended for situations such as:
//   • Interacting with APIs that use pointer fields to express optional values (e.g. json, yaml).
//   • Writing concise tests that need pointer literals.
//   • Reducing repetitive nil checks when consuming optional data.
//
// Usage examples can be found in example_test.go within this package.
