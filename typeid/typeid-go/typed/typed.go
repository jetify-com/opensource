package typed

import (
	"fmt"

	untyped "go.jetpack.io/typeid"
)

// IDType is an interface used to represent a statically checked ID type.
// Example:
// type UserIDType struct{}
// func (UserIDType) Prefix() string { return "user" }
type IDType interface {
	Type() string
}

// TypeID is a unique identifier with a given type as defined by the TypeID spec
type TypeID[T IDType] untyped.TypeID

// New returns a new TypeID with a random suffix and the given type.
func New[T IDType]() (TypeID[T], error) {
	var prefix T
	tid, err := untyped.New(prefix.Type())
	if err != nil {
		// Clients should ignore the id value when an error is present, but just
		// in case, construct a "nil" id of the given type.
		return Nil[T](), err
	}
	return (TypeID[T])(tid), nil
}

// Nil returns the null typeid of the given type.
func Nil[T IDType]() TypeID[T] {
	var prefix T
	return TypeID[T](untyped.Must(untyped.From(prefix.Type(), "00000000000000000000000000")))
}

// Type returns the type prefix of the TypeID
func (tid TypeID[T]) Type() string {
	return untyped.TypeID(tid).Type()
}

// Suffix returns the suffix of the TypeID in it's canonical base32 representation.
func (tid TypeID[T]) Suffix() string {
	return untyped.TypeID(tid).Suffix()
}

// String returns the TypeID in it's canonical string representation of the form:
// <prefix>_<suffix> where <suffix> is the canonical base32 representation of the UUID
func (tid TypeID[T]) String() string {
	return untyped.TypeID(tid).String()
}

// UUIDBytes decodes the TypeID's suffix as a UUID and returns it's bytes
func (tid TypeID[T]) UUIDBytes() []byte {
	return untyped.TypeID(tid).UUIDBytes()
}

// UUID decode the TypeID's suffix as a UUID and returns it as a hex formatted string
func (tid TypeID[T]) UUID() string {
	return untyped.TypeID(tid).UUID()
}

// From returns a new TypeID of the given type using the provided suffix
func From[T IDType](suffix string) (TypeID[T], error) {
	var prefix T
	tid, err := untyped.From(prefix.Type(), suffix)
	if err != nil {
		return Nil[T](), err
	}
	return (TypeID[T])(tid), nil
}

// FromString parses a TypeID from the given string. Returns an error if the
// string is not a valid TypeID, OR if the type prefix does not match the
// expected type.
func FromString[T IDType](s string) (TypeID[T], error) {
	var prefix T
	tid, err := untyped.FromString(s)
	if err != nil {
		return Nil[T](), err
	}
	if tid.Type() != prefix.Type() {
		return Nil[T](), fmt.Errorf("invalid type, expected %s but got %s", prefix.Type(), tid.Type())
	}
	return (TypeID[T])(tid), nil
}

// FromUUID returns a new TypeID of the given type using the provided UUID
func FromUUID[T IDType](uuid string) (TypeID[T], error) {
	var prefix T
	tid, err := untyped.FromUUID(prefix.Type(), uuid)
	if err != nil {
		return Nil[T](), err
	}
	return (TypeID[T])(tid), nil
}

// FromUUIDBytes returns a new TypeID of the given type using the provided UUID bytes
func FromUUIDBytes[T IDType](uuid []byte) (TypeID[T], error) {
	var prefix T
	tid, err := untyped.FromUUIDBytes(prefix.Type(), uuid)
	if err != nil {
		return Nil[T](), err
	}
	return (TypeID[T])(tid), nil
}

// Must panics if the given error is non-nil, otherwise it returns the given TypeID
func Must[T IDType](tid TypeID[T], err error) TypeID[T] {
	if err != nil {
		panic(err)
	}
	return tid
}
