package typeid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"go.jetpack.io/typeid/base32"
)

// New returns a new TypeID of the given type with a random suffix.
func New[T Subtype, PT subtypePtr[T]]() (T, error) {
	var nilID T

	prefix := subtypePrefix[T]()
	if prefix == anyPrefix {
		return nilID, errors.New("constructor error: use WithPrefix(), New() is for Subtypes")
	}
	return from[T, PT](prefix, "")
}

// WithPrefix returns a new TypeID with the given prefix and a random suffix.
// If you want to create an id without a prefix, pass an empty string.
func WithPrefix(prefix string) (TypeID, error) {
	return from[TypeID](prefix, "")
}

// From returns a new TypeID with the given prefix and suffix.
// If suffix is the empty string, a random suffix will be generated.
// If you want to create an id without a prefix, pass an empty string as the prefix.
func From(prefix string, suffix string) (TypeID, error) {
	return from[TypeID](prefix, suffix)
}

func FromSuffix[T Subtype, PT subtypePtr[T]](suffix string) (T, error) {
	var nilID T

	prefix := subtypePrefix[T]()
	if prefix == anyPrefix {
		return nilID, errors.New("constructor error: use From(prefix, suffix), FromSuffix is for Subtypes")
	}
	return from[T, PT](prefix, suffix)
}

// FromString parses a TypeID from a string of the form <prefix>_<suffix>
func FromString(s string) (TypeID, error) {
	return Parse[TypeID](s)
}

// Parse parses a TypeID from a string of the form <prefix>_<suffix>
// and ensures the TypeID is of the right type.
func Parse[T Subtype, PT subtypePtr[T]](s string) (T, error) {
	var nilID T
	prefix, suffix, err := split(s)
	if err != nil {
		return nilID, err
	}
	return from[T, PT](prefix, suffix)
}

func split(id string) (string, string, error) {
	switch parts := strings.SplitN(id, "_", 2); len(parts) {
	case 1:
		return "", parts[0], nil
	case 2:
		if parts[0] == "" {
			return "", "", errors.New("prefix cannot be empty when there's a separator")
		}
		return parts[0], parts[1], nil
	default:
		return "", "", fmt.Errorf("invalid typeid: %s", id)
	}
}

// FromUUID encodes the given UUID (in hex string form) as a TypeID with the given prefix.
func FromUUID[T Subtype, PT subtypePtr[T]](prefix string, uidStr string) (T, error) {
	uid, err := uuid.FromString(uidStr)
	var nilID T

	if err != nil {
		return nilID, err
	}
	suffix := base32.Encode(uid)
	return from[T, PT](prefix, suffix)
}

// FromUUID encodes the given UUID (in byte form) as a TypeID with the given prefix.
func FromUUIDBytes[T Subtype, PT subtypePtr[T]](prefix string, bytes []byte) (T, error) {
	uidStr := uuid.FromBytesOrNil(bytes).String()
	return FromUUID[T, PT](prefix, uidStr)
}

func from[T Subtype, PT subtypePtr[T]](prefix string, suffix string) (T, error) {
	var nilID T

	if err := validatePrefix[T](prefix); err != nil {
		return nilID, err
	}

	if suffix == "" {
		uid, err := uuid.NewV7()
		if err != nil {
			return nilID, err
		}
		suffix = base32.Encode(uid)
	}

	if err := validateSuffix(suffix); err != nil {
		return nilID, err
	}

	result := newSubtype[T, PT](prefix, suffix)
	return result, nil
}
