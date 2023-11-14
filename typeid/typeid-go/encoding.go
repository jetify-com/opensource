package typeid

import (
	"encoding"
	"fmt"
)

// TODO: Define a standardized binary encoding for typeids in the spec
// and use that to implement encoding.BinaryMarshaler and encoding.BinaryUnmarshaler

var _ encoding.TextMarshaler = (*TypeID[AnyPrefix])(nil)
var _ encoding.TextUnmarshaler = (*TypeID[AnyPrefix])(nil)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses a TypeID from a string using the same logic as FromString()
func (tid *TypeID[P]) UnmarshalText(text []byte) error {
	parsed, err := Parse[TypeID[P]](string(text))
	if err != nil {
		return err
	}
	*tid = parsed
	return nil
}

// Scan implements the sql.Scanner interface
func (tid *TypeID[P]) Scan(src any) error {
	switch obj := src.(type) {
	case nil:
		return nil
	case string:
		return tid.UnmarshalText([]byte(obj))
	// TODO: add supporte for []byte
	// we don't just want to store the full string as a byte array. Instead
	// we should encode using the UUID bytes. We could add support for
	// Binary Marshalling and Unmarshalling at the same time.
	default:
		return fmt.Errorf("unsupported scan type %T", obj)
	}
}

// MarshalText implements the encoding.TextMarshaler interface.
// It encodes a TypeID as a string using the same logic as String()
func (tid TypeID[P]) MarshalText() (text []byte, err error) {
	encoded := tid.String()
	return []byte(encoded), nil
}
