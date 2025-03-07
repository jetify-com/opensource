package try_test

import (
	"errors"
	"fmt"

	"go.jetify.com/pkg/try"
)

// This example shows how to create successful Try values
func ExampleOk() {
	// Create a successful result
	r := try.Ok(42)
	fmt.Println("Is ok:", r.IsOk())
	fmt.Println("Value:", r.MustGet())
	// Output:
	// Is ok: true
	// Value: 42
}

// This example shows how to wrap a value-error pair into a Try
func ExampleWrap() {
	// Function that returns a value and an error
	someFunction := func() (int, error) {
		return 42, nil
	}

	// Wrap the result
	r := try.Wrap(someFunction())
	fmt.Println("Is ok:", r.IsOk())
	fmt.Println("Value:", r.MustGet())
	// Output:
	// Is ok: true
	// Value: 42
}

// This example shows how to create error results
func ExampleErr() {
	// Create error results
	r1 := try.Err[int](errors.New("something went wrong"))
	r2 := try.Errf[int]("failed: %v", "bad input")

	fmt.Println("r1 is ok:", r1.IsOk())
	fmt.Println("r1 error:", r1.Err())
	fmt.Println("r2 is ok:", r2.IsOk())
	fmt.Println("r2 error:", r2.Err())
	// Output:
	// r1 is ok: false
	// r1 error: something went wrong
	// r2 is ok: false
	// r2 error: failed: bad input
}

// This example shows how to check the result state
func ExampleTry_IsOk() {
	r := try.Ok(42)

	if r.IsOk() {
		value := r.MustGet()
		fmt.Println("Got value:", value)
	} else {
		fmt.Println("Got error:", r.Err())
	}
	// Output:
	// Got value: 42
}

// This example shows safe error handling with Unwrap
func ExampleTry_Unwrap() {
	// Successful result
	r1 := try.Ok(42)
	if value, err := r1.Unwrap(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value)
	}

	// Error result
	r2 := try.Err[int](errors.New("something went wrong"))
	if value, err := r2.Unwrap(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value)
	}
	// Output:
	// Value: 42
	// Error: something went wrong
}

// This example shows how to provide fallback values
func ExampleTry_OrElse() {
	// Successful result
	r1 := try.Ok(42)
	v1 := r1.OrElse(0)
	fmt.Println("Value with fallback:", v1)

	// Error result
	r2 := try.Err[int](errors.New("something went wrong"))
	v2 := r2.OrElse(0)
	fmt.Println("Value with fallback for error:", v2)
	// Output:
	// Value with fallback: 42
	// Value with fallback for error: 0
}

// This example shows how to convert panics to Results
func ExampleDo() {
	// A function that might panic
	riskyOperation := func() int {
		// Uncomment to see panic handling
		// panic("something went terribly wrong")
		return 42
	}

	// Convert potential panic to Try
	r := try.Do(func() int {
		return riskyOperation()
	})

	if r.IsOk() {
		fmt.Println("Value:", r.MustGet())
	} else {
		fmt.Println("Error:", r.Err())
	}
	// Output:
	// Value: 42
}

// This example shows how to handle async operations
func ExampleGo() {
	// Define a function that simulates a complex calculation
	complexCalculation := func() (int, error) {
		// Simulate work
		return 42, nil
	}

	// Run the calculation asynchronously
	ch := try.Go(func() (int, error) {
		return complexCalculation()
	})

	// Receive the result when ready
	r := <-ch

	if r.IsOk() {
		fmt.Println("Async result:", r.MustGet())
	} else {
		fmt.Println("Async error:", r.Err())
	}
	// Output:
	// Async result: 42
}
