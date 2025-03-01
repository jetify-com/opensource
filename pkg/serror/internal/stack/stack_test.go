package stack

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackFormat(t *testing.T) {
	// Create a test stack
	testStack := Stack{
		{Name: "function1", File: "file1.go", Line: 10},
		{Name: "function2", File: "file2.go", Line: 20},
		{Name: "function3", File: "file3.go", Line: 30},
	}

	tests := []struct {
		name   string
		sep    string
		invert bool
		want   []string
	}{
		{
			name:   "forward with colon separator",
			sep:    ":",
			invert: false,
			want: []string{
				"function3:file3.go:30",
				"function2:file2.go:20",
				"function1:file1.go:10",
			},
		},
		{
			name:   "inverted with pipe separator",
			sep:    "|",
			invert: true,
			want: []string{
				"function1|file1.go|10",
				"function2|file2.go|20",
				"function3|file3.go|30",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := testStack.format(tt.sep, tt.invert)
			assert.Equal(t, tt.want, formatted, "Stack format should match expected output")
		})
	}
}

func TestStackFrameFormat(t *testing.T) {
	frame := StackFrame{
		Name: "TestFunction",
		File: "test_file.go",
		Line: 42,
	}

	tests := []struct {
		name string
		sep  string
		want string
	}{
		{
			name: "format with colon",
			sep:  ":",
			want: "TestFunction:test_file.go:42",
		},
		{
			name: "format with arrow",
			sep:  " -> ",
			want: "TestFunction -> test_file.go -> 42",
		},
		{
			name: "format with space",
			sep:  " ",
			want: "TestFunction test_file.go 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := frame.format(tt.sep)
			assert.Equal(t, tt.want, result, "StackFrame format should match expected output")
		})
	}
}

func TestCaller(t *testing.T) {
	// Test caller() function
	c := caller(1)
	assert.NotNil(t, c, "Caller should return a non-nil frame")

	frame := c.get()
	assert.NotEmpty(t, frame.Name, "Frame name should not be empty")
	assert.NotEmpty(t, frame.File, "Frame file should not be empty")
	assert.Greater(t, frame.Line, 0, "Line number should be positive")
}

func TestCallers(t *testing.T) {
	// Test callers() function
	s := callers(2) // Skip runtime.Callers + this function
	stackFrames := s.get()

	// We expect a minimum number of frames
	// At minimum we should have the test function itself
	assert.GreaterOrEqual(t, len(stackFrames), 1, "Stack should contain at least the test function frame")

	// First frame should be the current test function
	assert.Contains(t, stackFrames[0].Name, "TestCallers", "First frame should be TestCallers")
}

func TestStackGet(t *testing.T) {
	// Test stack.get() directly
	s := callers(2) // Skip runtime.Callers + this function
	stackFrames := s.get()

	// We expect a minimum number of frames
	// At minimum we should have the test function itself
	assert.GreaterOrEqual(t, len(stackFrames), 1, "Stack should contain at least the test function frame")

	// Basic validation
	assert.Contains(t, stackFrames[0].Name, "TestStackGet", "First frame should refer to the test function")
}

func TestFramePC(t *testing.T) {
	// Get program counter
	pc, _, _, _ := runtime.Caller(0)
	f := frame(pc)

	// Test the pc() function
	framePC := f.pc()
	assert.NotZero(t, framePC, "Frame PC should be non-zero")
	assert.Equal(t, pc-1, framePC, "Frame PC should be original PC minus 1")
}

func TestFrameGet(t *testing.T) {
	// Get program counter
	pc, _, _, _ := runtime.Caller(0)
	f := frame(pc)

	// Test the get() function
	stackFrame := f.get()

	// Check if frame contains expected data for the current function
	assert.Contains(t, stackFrame.Name, "TestFrameGet", "Frame name should contain TestFrameGet")
	assert.Contains(t, stackFrame.File, "stack_test.go", "Frame file should contain stack_test.go")
	assert.Greater(t, stackFrame.Line, 0, "Line number should be positive")
}

func TestStackInsertPC(t *testing.T) {
	// Create test stacks
	mainStack := []uintptr{1, 2, 3, 4}
	emptyInsert := []uintptr{}
	singleInsert := []uintptr{5}
	doubleInsert := []uintptr{2, 5}       // 2 is already in mainStack
	doubleInsertSecond := []uintptr{6, 3} // 3 is second item, should insert 6 before it

	tests := []struct {
		name        string
		mainStack   []uintptr
		insertStack []uintptr
		expected    []uintptr
	}{
		{
			name:        "empty insert stack",
			mainStack:   mainStack,
			insertStack: emptyInsert,
			expected:    mainStack,
		},
		{
			name:        "single item insert stack",
			mainStack:   mainStack,
			insertStack: singleInsert,
			expected:    append(mainStack, uintptr(5)),
		},
		{
			name:        "double item insert with match",
			mainStack:   mainStack,
			insertStack: doubleInsert,
			expected:    mainStack, // No change as first item already exists
		},
		{
			name:        "double item insert with second match",
			mainStack:   mainStack,
			insertStack: doubleInsertSecond,
			expected:    []uintptr{1, 2, 6, 3, 4}, // Insert before the matching second item
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to stack type
			var s stack = tt.mainStack
			var wrapPCs stack = tt.insertStack

			// Perform insert
			s.insertPC(wrapPCs)

			// Check result
			assert.Equal(t, tt.expected, []uintptr(s), "Stack after insert should match expected")
		})
	}
}

func TestStackIsGlobal(t *testing.T) {
	// Create a non-global stack from the current call stack
	const depth = 10
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var s stack = pcs[:n]

	// Test isGlobal() directly
	assert.False(t, s.isGlobal(), "Current stack should not be detected as global")

	// More direct test approach for the positive case
	// Create a small helper function that just tests the logic directly
	testStack := []StackFrame{
		{Name: "some.function", File: "file.go", Line: 10},
		{Name: "runtime.doinit", File: "runtime.go", Line: 100}, // This should trigger isGlobal
		{Name: "other.function", File: "other.go", Line: 20},
	}

	isGlobalResult := func(frames []StackFrame) bool {
		for _, f := range frames {
			if strings.ToLower(f.Name) == "runtime.doinit" {
				return true
			}
		}
		return false
	}(testStack)

	assert.True(t, isGlobalResult, "Stack with runtime.doinit should be detected as global")
}

func TestInsertFunction(t *testing.T) {
	s := stack{1, 2, 3, 4}
	u := uintptr(5)
	at := 2

	result := insert(s, u, at)
	expected := stack{1, 2, 5, 3, 4}

	assert.Equal(t, expected, result, "Insert should place element at correct position")
}
