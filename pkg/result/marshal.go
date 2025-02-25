package result

import (
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v3"
)

type encodedResult[T any] struct {
	Value T      `json:"value,omitempty" yaml:"value,omitempty" toml:"value,omitempty"`
	Error string `json:"error,omitempty" yaml:"error,omitempty" toml:"error,omitempty"`
}

// toEncoding converts a Result to its encoding representation
func (r Result[T]) toEncoding() encodedResult[T] {
	if r.IsErr() {
		return encodedResult[T]{
			Error: r.err.Error(),
		}
	}
	return encodedResult[T]{
		Value: r.value,
	}
}

// fromEncoding converts from an encoding representation to a Result
func fromEncoding[T any](enc encodedResult[T]) Result[T] {
	if enc.Error != "" {
		return Err[T](errors.New(enc.Error))
	}
	return Ok(enc.Value)
}

// JSON Marshaling
// --------------

// MarshalJSON serializes the Result into JSON. For successful results,
// it produces {"value": ...}; for errors, it produces {"error": "..."}.
func (r Result[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.toEncoding())
}

// UnmarshalJSON deserializes from JSON into a Result. It expects either
// {"value": ...} or {"error": "..."}.
func (r *Result[T]) UnmarshalJSON(data []byte) error {
	var enc encodedResult[T]
	if err := json.Unmarshal(data, &enc); err != nil {
		return err
	}
	*r = fromEncoding(enc)
	return nil
}

// YAML Marshaling
// --------------

// MarshalYAML serializes the Result into YAML.
func (r Result[T]) MarshalYAML() (any, error) {
	return r.toEncoding(), nil
}

// UnmarshalYAML deserializes from YAML into a Result.
func (r *Result[T]) UnmarshalYAML(value *yaml.Node) error {
	var enc encodedResult[T]
	if err := value.Decode(&enc); err != nil {
		return err
	}
	*r = fromEncoding(enc)
	return nil
}
