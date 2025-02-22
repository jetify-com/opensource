// Package result provides a generic Result type for handling operations that
// can either succeed with a value or fail with an error.
//
// The Result type is particularly useful when you want to:
//   - Defer error handling
//   - Chain operations that might fail
//   - Handle both synchronous and asynchronous operations that can fail
//   - Convert panics into errors
//   - Process a list of tasks, each of which might fail
//
// Basic usage:
//
//	// Create successful results
//	r1 := result.Ok(42)
//	r2 := result.From(someFunction()) // from (value, error) pair
//
//	// Create error results
//	r3 := result.Err[int](errors.New("something went wrong"))
//	r4 := result.Errf[int]("failed: %v", err)
//
//	// Check result state
//	if r1.IsOk() {
//	    value := r1.MustGet()
//	    // use value...
//	}
//
//	// Safe error handling
//	if value, err := r1.Get(); err != nil {
//	    // handle error...
//	} else {
//	    // use value...
//	}
//
//	// Provide fallback values
//	value := r1.OrElse(0) // returns 0 if r1 contains an error
//
// The package also provides utilities for handling panics and asynchronous
// operations:
//
//	// Convert panics to Results
//	r := result.Do(func() int {
//	    // this will be caught if it panics
//	    return riskyOperation()
//	})
//
//	// Handle async operations
//	ch := result.Go(func() (int, error) {
//	    // async operation
//	    return complexCalculation()
//	})
//	r := <-ch // receive Result when ready
package result
