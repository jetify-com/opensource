package try

import (
	"encoding/json"
	"errors"
)

type encodedTry[T any] struct {
	Value T      `json:"value,omitempty" yaml:"value,omitempty" toml:"value,omitempty"`
	Error string `json:"error,omitempty" yaml:"error,omitempty" toml:"error,omitempty"`
}

// toEncoding converts a Try to its encoding representation
func (r Try[T]) toEncoding() encodedTry[T] {
	if r.IsErr() {
		return encodedTry[T]{
			Error: r.err.Error(),
		}
	}
	return encodedTry[T]{
		Value: r.value,
	}
}

// fromEncoding converts from an encoding representation to a Try
func fromEncoding[T any](enc encodedTry[T]) Try[T] {
	if enc.Error != "" {
		return Err[T](errors.New(enc.Error))
	}
	return Ok(enc.Value)
}

// JSON Marshaling
// --------------

// MarshalJSON serializes the Try into JSON. For successful results,
// it produces {"value": ...}; for errors, it produces {"error": "..."}.
func (r Try[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.toEncoding())
}

// UnmarshalJSON deserializes from JSON into a Try. It expects either
// {"value": ...} or {"error": "..."}.
func (r *Try[T]) UnmarshalJSON(data []byte) error {
	var enc encodedTry[T]
	if err := json.Unmarshal(data, &enc); err != nil {
		return err
	}
	*r = fromEncoding(enc)
	return nil
}

// YAML Marshaling
// --------------

// MarshalYAML serializes the Try into YAML.
func (r Try[T]) MarshalYAML() (any, error) {
	return r.toEncoding(), nil
}

// UnmarshalYAML deserializes from YAML into a Try. It expects either
// value: ... or error: ...
func (r *Try[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var enc encodedTry[T]
	if err := unmarshal(&enc); err != nil {
		return err
	}
	*r = fromEncoding(enc)
	return nil
}
