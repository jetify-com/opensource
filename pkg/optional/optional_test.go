package optional

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSome(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected Option[int]
	}{
		{
			name:     "creates option with integer value",
			value:    42,
			expected: Option[int]{value: 42, set: true},
		},
		{
			name:     "creates option with zero value",
			value:    0,
			expected: Option[int]{value: 0, set: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Some(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNone(t *testing.T) {
	tests := []struct {
		name     string
		expected Option[int]
	}{
		{
			name:     "creates empty option",
			expected: Option[int]{set: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := None[int]()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name          string
		option        Option[string]
		expectedValue string
		expectedSet   bool
	}{
		{
			name:          "gets value from Some",
			option:        Some("hello"),
			expectedValue: "hello",
			expectedSet:   true,
		},
		{
			name:          "gets zero value from None",
			option:        None[string](),
			expectedValue: "",
			expectedSet:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, set := tt.option.Get()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedSet, set)
		})
	}
}

func TestMustGet(t *testing.T) {
	tests := []struct {
		name          string
		option        Option[float64]
		expectedValue float64
		shouldPanic   bool
	}{
		{
			name:          "gets value from Some",
			option:        Some(3.14),
			expectedValue: 3.14,
			shouldPanic:   false,
		},
		{
			name:        "panics when accessing None",
			option:      None[float64](),
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					tt.option.MustGet()
				})
			} else {
				value := tt.option.MustGet()
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}

func TestOrElse(t *testing.T) {
	tests := []struct {
		name          string
		option        Option[int]
		defaultValue  int
		expectedValue int
	}{
		{
			name:          "returns contained value for Some",
			option:        Some(42),
			defaultValue:  0,
			expectedValue: 42,
		},
		{
			name:          "returns default value for None",
			option:        None[int](),
			defaultValue:  99,
			expectedValue: 99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.option.GetOrElse(tt.defaultValue)
			assert.Equal(t, tt.expectedValue, result)
		})
	}
}

func TestIsNone(t *testing.T) {
	tests := []struct {
		name     string
		option   Option[string]
		expected bool
	}{
		{
			name:     "returns false for Some",
			option:   Some("value"),
			expected: false,
		},
		{
			name:     "returns true for None",
			option:   None[string](),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.option.IsNone()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		option   Option[int]
		expected string
	}{
		{
			name:     "formats Some value",
			option:   Some(42),
			expected: "Some(42)",
		},
		{
			name:     "formats None",
			option:   None[int](),
			expected: "None",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.option.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGoString(t *testing.T) {
	tests := []struct {
		name     string
		option   Option[any]
		expected string
	}{
		{
			name:     "formats Some value with int",
			option:   Some[any](42),
			expected: "Some[int](42)",
		},
		{
			name:     "formats Some value with string",
			option:   Some[any]("test"),
			expected: `Some[string]("test")`,
		},
		{
			name:     "formats None",
			option:   None[any](),
			expected: "None[<nil>]()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.option.GoString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		option   Option[any]
		expected string
	}{
		{
			name:     "marshals Some integer",
			option:   Some[any](42),
			expected: "42",
		},
		{
			name:     "marshals Some string",
			option:   Some[any]("test"),
			expected: `"test"`,
		},
		{
			name:     "marshals Some object",
			option:   Some[any](map[string]int{"a": 1}),
			expected: `{"a":1}`,
		},
		{
			name:     "marshals None as null",
			option:   None[any](),
			expected: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.option)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected Option[any]
	}{
		{
			name:     "unmarshals integer into Some",
			json:     "42",
			expected: Some[any](float64(42)), // JSON numbers unmarshal as float64 by default
		},
		{
			name:     "unmarshals string into Some",
			json:     `"test"`,
			expected: Some[any]("test"),
		},
		{
			name:     "unmarshals object into Some",
			json:     `{"a":1}`,
			expected: Some[any](map[string]interface{}{"a": float64(1)}),
		},
		{
			name:     "unmarshals null into None",
			json:     "null",
			expected: None[any](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Option[any]
			err := json.Unmarshal([]byte(tt.json), &result)
			require.NoError(t, err)

			// Check if None or Some
			assert.Equal(t, tt.expected.IsNone(), result.IsNone())

			if !tt.expected.IsNone() {
				expectedValue, _ := tt.expected.Get()
				actualValue, _ := result.Get()
				assert.Equal(t, expectedValue, actualValue)
			}
		})
	}
}

// TestUnmarshalJSONEdgeCases tests edge cases for the UnmarshalJSON implementation
func TestUnmarshalJSONEdgeCases(t *testing.T) {
	t.Run("invalid json", func(t *testing.T) {
		var result Option[int]
		err := json.Unmarshal([]byte("{invalid json}"), &result)
		assert.Error(t, err)
	})

	t.Run("type mismatch", func(t *testing.T) {
		var result Option[int]
		err := json.Unmarshal([]byte(`"not an int"`), &result)
		assert.Error(t, err)
	})

	t.Run("custom json marshaler type", func(t *testing.T) {
		type CustomType struct {
			Value string
		}

		// Create a nested option with a custom type
		original := Some(Some(CustomType{Value: "test"}))
		data, err := json.Marshal(original)
		require.NoError(t, err)

		var result Option[Option[CustomType]]
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Verify the nested structure was preserved
		val1, set1 := result.Get()
		assert.True(t, set1)

		val2, set2 := val1.Get()
		assert.True(t, set2)
		assert.Equal(t, "test", val2.Value)
	})
}

// TestStructWithOptions tests using Option with struct fields
func TestStructWithOptions(t *testing.T) {
	type Person struct {
		Name    string         `json:"name"`
		Age     Option[int]    `json:"age,omitzero"` // Will be omitted when None
		Address Option[string] `json:"address"`      // Will be null when None
	}

	tests := []struct {
		name     string
		person   Person
		expected string
	}{
		{
			name: "all fields present",
			person: Person{
				Name:    "Alice",
				Age:     Some(30),
				Address: Some("123 Main St"),
			},
			expected: `{"name":"Alice","age":30,"address":"123 Main St"}`,
		},
		{
			name: "all optional fields as None",
			person: Person{
				Name:    "Bob",
				Age:     None[int](),
				Address: None[string](),
			},
			expected: `{"name":"Bob","address":null}`, // Age omitted, Address null
		},
		{
			name: "mixed Some and None",
			person: Person{
				Name:    "Charlie",
				Age:     Some(25),
				Address: None[string](),
			},
			expected: `{"name":"Charlie","age":25,"address":null}`, // Address still null
		},
		{
			name: "omitted field but present null field",
			person: Person{
				Name:    "Dave",
				Age:     None[int](),
				Address: Some("456 Oak St"),
			},
			expected: `{"name":"Dave","address":"456 Oak St"}`, // Age omitted
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.person)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))

			// Test unmarshaling
			var result Person
			err = json.Unmarshal(data, &result)
			require.NoError(t, err)

			assert.Equal(t, tt.person.Name, result.Name)

			origAge, origAgeSet := tt.person.Age.Get()
			resultAge, resultAgeSet := result.Age.Get()
			assert.Equal(t, origAgeSet, resultAgeSet)
			if origAgeSet {
				assert.Equal(t, origAge, resultAge)
			}

			origAddr, origAddrSet := tt.person.Address.Get()
			resultAddr, resultAddrSet := result.Address.Get()
			assert.Equal(t, origAddrSet, resultAddrSet)
			if origAddrSet {
				assert.Equal(t, origAddr, resultAddr)
			}
		})
	}
}

// TestPtr tests the Ptr function for creating Option values from pointers
func TestPtr(t *testing.T) {
	tests := []struct {
		name     string
		ptr      *string
		expected Option[string]
	}{
		{
			name:     "creates Some from non-nil pointer",
			ptr:      stringPtr("test"),
			expected: Some("test"),
		},
		{
			name:     "creates None from nil pointer",
			ptr:      nil,
			expected: None[string](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Ptr(tt.ptr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// TestZeroValueVsNone tests the distinction between a Some with a zero value and None
func TestZeroValueVsNone(t *testing.T) {
	tests := []struct {
		name          string
		option        Option[any]
		isNone        bool
		expectedValue any
	}{
		{
			name:          "Some with zero int",
			option:        Some[any](0),
			isNone:        false,
			expectedValue: 0,
		},
		{
			name:          "Some with empty string",
			option:        Some[any](""),
			isNone:        false,
			expectedValue: "",
		},
		{
			name:          "Some with nil",
			option:        Some[any](nil),
			isNone:        false,
			expectedValue: nil,
		},
		{
			name:          "None",
			option:        None[any](),
			isNone:        true,
			expectedValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test IsNone distinguishes None from Some with zero values
			assert.Equal(t, tt.isNone, tt.option.IsNone())

			// For Some values, verify the value is as expected
			if !tt.isNone {
				value, set := tt.option.Get()
				assert.True(t, set)
				assert.Equal(t, tt.expectedValue, value)
			}

			// Verify JSON marshaling works correctly
			data, err := json.Marshal(tt.option)
			require.NoError(t, err)

			if tt.isNone {
				assert.JSONEq(t, "null", string(data))
			} else if tt.expectedValue == nil {
				assert.JSONEq(t, "null", string(data))
			} else {
				expectedJSON, err := json.Marshal(tt.expectedValue)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedJSON), string(data))
			}
		})
	}
}
