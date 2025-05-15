package pointer_test

import (
	"fmt"

	"go.jetify.com/pkg/pointer"
)

// Example demonstrates the basic pointer creation functions.
func Example() {
	// Create pointers to primitive values
	boolPtr := pointer.Bool(true)
	intPtr := pointer.Int(42)
	stringPtr := pointer.String("hello")

	// Safely dereference the pointers
	fmt.Printf("Bool: %v\n", *boolPtr)
	fmt.Printf("Int: %v\n", *intPtr)
	fmt.Printf("String: %v\n", *stringPtr)

	// Output:
	// Bool: true
	// Int: 42
	// String: hello
}

// ExamplePtr demonstrates the use of the generic Ptr function.
func ExamplePtr() {
	// Create pointers to primitive values
	boolPtr := pointer.Ptr(true)
	intPtr := pointer.Ptr(42)
	stringPtr := pointer.Ptr("hello")

	// Safely dereference the pointers
	fmt.Printf("Bool: %v\n", *boolPtr)
	fmt.Printf("Int: %v\n", *intPtr)
	fmt.Printf("String: %v\n", *stringPtr)

	// Output:
	// Bool: true
	// Int: 42
	// String: hello
}

// ExampleValueOr demonstrates the use of the ValueOr function with nil pointers.
func ExampleValueOr() {
	var nilStringPtr *string
	definedStringPtr := pointer.String("defined value")

	// Use ValueOr to handle nil pointers with defaults
	nilValue := pointer.ValueOr(nilStringPtr, "default value")
	definedValue := pointer.ValueOr(definedStringPtr, "default value")

	fmt.Printf("Nil pointer value: %s\n", nilValue)
	fmt.Printf("Defined pointer value: %s\n", definedValue)

	// Output:
	// Nil pointer value: default value
	// Defined pointer value: defined value
}

// ExampleValue demonstrates the use of the Value function.
func ExampleValue() {
	var nilIntPtr *int
	definedIntPtr := pointer.Int(42)

	// Use Value to safely dereference pointers
	nilValue := pointer.Value(nilIntPtr)
	definedValue := pointer.Value(definedIntPtr)

	fmt.Printf("Nil pointer value: %d\n", nilValue)
	fmt.Printf("Defined pointer value: %d\n", definedValue)

	// Output:
	// Nil pointer value: 0
	// Defined pointer value: 42
}
