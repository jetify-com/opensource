package typed

import (
	"fmt"

	untyped "go.jetpack.io/typeid"
)

// Create a typed typeid by creating a new string type:
//
// type userID string
// const userPrefix = userID("user")
//
// typeid.T(userPrefix).New()
// typeid.T(userPrefix).FromString("user-1234")

type TypeID[T ~string] struct {
	untyped.TypeID
}

// This can be removed once untyped.TypeID can handle nil values
func (t *TypeID[T]) String() string {
	if t == nil {
		return untyped.Nil.String()
	}
	return t.TypeID.String()
}

// This can be removed once untyped.TypeID can handle nil values
func (t *TypeID[T]) Type() string {
	if t == nil {
		return ""
	}
	return t.TypeID.Type()
}

type typeIDBuilder[T ~string] struct {
	prefix string
}

func T[T ~string](prefix T) *typeIDBuilder[T] {
	return &typeIDBuilder[T]{string(prefix)}
}

// New returns a new TypeID with a random suffix and the given type.
func (b *typeIDBuilder[T]) New() (*TypeID[T], error) {
	tid, err := untyped.New(b.prefix)
	if err != nil {
		// Clients should ignore the id value when an error is present, but just
		// in case, construct a "nil" id of the given type.
		return nil, err
	}
	return &TypeID[T]{tid}, nil
}

// From returns a new TypeID of the given type using the provided suffix
func (b *typeIDBuilder[T]) From(suffix string) (*TypeID[T], error) {
	id, err := untyped.From(b.prefix, suffix)
	if err != nil {
		return nil, err
	}
	return &TypeID[T]{id}, nil
}

// FromString parses a TypeID from the given string. Returns an error if the
// string is not a valid TypeID, OR if the type prefix does not match the
// expected type.
func (b *typeIDBuilder[T]) FromString(s string) (*TypeID[T], error) {
	tid, err := untyped.FromString(s)
	if err != nil {
		return nil, err
	}
	if tid.Type() != b.prefix {
		return nil, fmt.Errorf("invalid type, expected %s but got %s", b.prefix, tid.Type())
	}
	return &TypeID[T]{tid}, nil
}

// FromUUID returns a new TypeID of the given type using the provided UUID
func (b *typeIDBuilder[T]) FromUUID(uuid string) (*TypeID[T], error) {
	tid, err := untyped.FromUUID(b.prefix, uuid)
	if err != nil {
		return nil, err
	}
	return &TypeID[T]{tid}, nil
}

// FromUUIDBytes returns a new TypeID of the given type using the provided UUID bytes
func (b *typeIDBuilder[T]) FromUUIDBytes(uuid []byte) (*TypeID[T], error) {
	tid, err := untyped.FromUUIDBytes(b.prefix, uuid)
	if err != nil {
		return nil, err
	}
	return &TypeID[T]{tid}, nil
}

// Must panics if the given error is non-nil, otherwise it returns the given TypeID
func (b *typeIDBuilder[T]) Must(tid TypeID[T], err error) TypeID[T] {
	if err != nil {
		panic(err)
	}
	return tid
}
