package typeid

import (
	"github.com/gofrs/uuid/v5"
	"go.jetpack.io/typeid/base32"
)

// TypeID is a unique identifier with a given type as defined by the TypeID spec
type TypeID struct {
	prefix string
	suffix string
}

// Nil represents the null TypeID: a type id with no prefix, and where the suffix
// is all zeroes.
var Nil = TypeID{
	prefix: "",
	suffix: "", // Will become "00000000000000000000000000" in string representation
}

const nilSuffix = "00000000000000000000000000"

// Type returns the type prefix of the TypeID
// Deprecated: use Prefix() instead.
func (tid TypeID) Type() string {
	return tid.Prefix()
}

// Prefix returns the type prefix of the TypeID
func (tid TypeID) Prefix() string {
	// If we're dealing with a subtype, use the subtype's hardcoded prefix:
	if tid.prefix == "" && tid.AllowedPrefix() != anyPrefix {
		return tid.AllowedPrefix()
	}
	return tid.prefix
}

const anyPrefix = "*"

func (tid TypeID) AllowedPrefix() string {
	return anyPrefix
}

// Suffix returns the suffix of the TypeID in it's canonical base32 representation.
func (tid TypeID) Suffix() string {
	// We want to treat the "empty" TypeID as equivalent to the Nil typeid
	if tid.suffix == "" {
		return nilSuffix
	}
	return tid.suffix
}

// String returns the TypeID in it's canonical string representation of the form:
// <prefix>_<suffix> where <suffix> is the canonical base32 representation of the UUID
func (tid TypeID) String() string {
	if tid.Prefix() == "" {
		return tid.Suffix()
	}
	return tid.Prefix() + "_" + tid.Suffix()
}

// UUIDBytes decodes the TypeID's suffix as a UUID and returns it's bytes
func (tid TypeID) UUIDBytes() []byte {
	b, err := base32.Decode(tid.Suffix())

	// Decode only fails if the suffix cannot be decoded for one of two reasons:
	// 1. The suffix is not 26 characters long
	// 2. The suffix contains characters that are not in the base32 alphabet
	// We guarantee that the suffix is valid in the TypeID constructors, so this panic
	// should never be reached.
	if err != nil {
		panic(err)
	}
	return b
}

// UUID decodes the TypeID's suffix as a UUID and returns it as a hex string
func (tid TypeID) UUID() string {
	return uuid.FromBytesOrNil(tid.UUIDBytes()).String()
}

// Must returns a TypeID if the error is nil, otherwise panics.
// Often used with New() to create a TypeID in a single line as follows:
// tid := Must(New("prefix"))
func Must[T Subtype](tid T, err error) T {
	if err != nil {
		panic(err)
	}
	return tid
}
