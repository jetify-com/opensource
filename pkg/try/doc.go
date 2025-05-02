// Package try provides a generic Try type for handling operations that
// can either succeed with a value or fail with an error.
//
// The Try type is particularly useful when you want to:
//   - Defer error handling
//   - Chain operations that might fail
//   - Handle both synchronous and asynchronous operations that can fail
//   - Process a list of tasks, each of which might fail
//
// Basic usage:
//
//	// Create successful results
//	r1 := try.Ok(42)
//	r2 := try.From(someFunction()) // from (value, error) pair
//
//	// Create error results
//	r3 := try.Err[int](errors.New("something went wrong"))
//	r4 := try.Errf[int]("failed: %v", err)
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
//	value := r1.GetOrElse(0) // returns 0 if r1 contains an error
package try
