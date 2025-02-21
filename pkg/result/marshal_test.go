package result

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type encoder interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

func TestMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		input    Result[string]
		expected map[string]string // expected output for each encoder
	}{
		{
			name:  "success value",
			input: Ok("test value"),
			expected: map[string]string{
				"json": `{"value":"test value"}`,
				"yaml": "value: test value\n",
			},
		},
		{
			name:  "error value",
			input: Err[string](errors.New("test error")),
			expected: map[string]string{
				"json": `{"error":"test error"}`,
				"yaml": "error: test error\n",
			},
		},
	}

	encoders := map[string]encoder{
		"json": &jsonEncoder{},
		"yaml": &yamlEncoder{},
	}

	for _, tt := range tests {
		for encName, enc := range encoders {
			t.Run(tt.name+"_"+encName, func(t *testing.T) {
				// Test marshaling
				data, err := enc.Marshal(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected[encName], string(data))

				// Test unmarshaling
				var actual Result[string]
				err = enc.Unmarshal([]byte(tt.expected[encName]), &actual)
				assert.NoError(t, err)

				if tt.input.IsErr() {
					assert.True(t, actual.IsErr())
					assert.Equal(t, tt.input.Err(), actual.Err())
				} else {
					assert.False(t, actual.IsErr())
					assert.Equal(t, tt.input.value, actual.value)
				}
			})
		}
	}

	// Test invalid input for each encoder
	invalidTests := map[string]string{
		"json": `{"invalid`,
		"yaml": `invalid: - yaml: content`,
	}

	for encName, enc := range encoders {
		t.Run("invalid_"+encName, func(t *testing.T) {
			var actual Result[string]
			err := enc.Unmarshal([]byte(invalidTests[encName]), &actual)
			assert.Error(t, err)
		})
	}
}

func TestComplexTypeMarshaling(t *testing.T) {
	type Complex struct {
		Name    string `json:"name" yaml:"name"`
		Numbers []int  `json:"numbers" yaml:"numbers"`
		Nested  struct {
			Value bool `json:"value" yaml:"value"`
		} `json:"nested" yaml:"nested"`
	}

	complex := Complex{
		Name:    "test",
		Numbers: []int{1, 2, 3},
		Nested: struct {
			Value bool `json:"value" yaml:"value"`
		}{
			Value: true,
		},
	}

	input := Ok(complex)

	encoders := map[string]encoder{
		"json": &jsonEncoder{},
		"yaml": &yamlEncoder{},
	}

	for name, enc := range encoders {
		t.Run("complex "+name, func(t *testing.T) {
			data, err := enc.Marshal(input)
			assert.NoError(t, err)

			var actual Result[Complex]
			err = enc.Unmarshal(data, &actual)
			assert.NoError(t, err)
			assert.Equal(t, input.value, actual.value)
		})
	}
}

func TestStructRoundtrip(t *testing.T) {
	type Person struct {
		Name    string   `json:"name" yaml:"name"`
		Age     int      `json:"age" yaml:"age"`
		Hobbies []string `json:"hobbies" yaml:"hobbies"`
	}

	tests := []struct {
		name  string
		input Result[Person]
	}{
		{
			name: "success case",
			input: Ok(Person{
				Name:    "Alice",
				Age:     30,
				Hobbies: []string{"reading", "coding"},
			}),
		},
		{
			name:  "error case",
			input: Err[Person](errors.New("person not found")),
		},
	}

	encoders := map[string]encoder{
		"json": &jsonEncoder{},
		"yaml": &yamlEncoder{},
	}

	for _, tt := range tests {
		for encName, enc := range encoders {
			t.Run(tt.name+"_"+encName, func(t *testing.T) {
				// Marshal the original value
				data, err := enc.Marshal(tt.input)
				require.NoError(t, err)

				// Unmarshal back into a new value
				var result Result[Person]
				err = enc.Unmarshal(data, &result)
				require.NoError(t, err)

				// Verify the roundtrip
				if tt.input.IsErr() {
					assert.True(t, result.IsErr())
					assert.Equal(t, tt.input.Err().Error(), result.Err().Error())
				} else {
					assert.False(t, result.IsErr())
					assert.Equal(t, tt.input.value, result.value)
				}
			})
		}
	}
}

// Encoder wrappers to implement common interface
type jsonEncoder struct{}

func (e *jsonEncoder) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (e *jsonEncoder) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

type yamlEncoder struct{}

func (e *yamlEncoder) Marshal(v any) ([]byte, error)      { return yaml.Marshal(v) }
func (e *yamlEncoder) Unmarshal(data []byte, v any) error { return yaml.Unmarshal(data, v) }
