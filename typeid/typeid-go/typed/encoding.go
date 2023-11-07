package typed

import "go.jetpack.io/typeid"

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses a TypeID from a string using the same logic as FromString()
func (tid *TypeID[T]) UnmarshalText(text []byte) error {
	tid2, err := typeid.FromString(string(text))
	if err != nil {
		return err
	}
	*tid = TypeID[T]{tid2}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It encodes a TypeID as a string using the same logic as String()
func (tid TypeID[T]) MarshalText() (text []byte, err error) {
	encoded := tid.String()
	return []byte(encoded), nil
}
