package typeid

// Subtype is an interface used to create a more specific subtype of TypeID
// For example, if you want to create an `OrgID` type that only accepts
// an `org_` prefix.
type Subtype interface {
	AllowedPrefix() string
	isTypeID() bool // Private to only allow implementation by embedding TypeID
}

var _ Subtype = (*TypeID)(nil)

type subtypePtr[T Subtype] interface {
	*T
	init(prefix string, suffix string)
}

func (tid *TypeID) init(prefix string, suffix string) {
	// In general TypeID is an immutable value-type, and pretty much every
	// "mutation" should return a copy with the modifications instead of modifying
	// the original. We make an exception for this *private* method, because
	// sometimes we need to modify the fields in the process of initializing
	// a new subtype.
	tid.prefix = prefix
	tid.suffix = suffix
}

func (tid TypeID) isTypeID() bool {
	return true
}

// Does not do any validation, use only as part of implementing a constructor
// that is doing the validation itself
func newSubtype[T Subtype, PT subtypePtr[T]](prefix string, suffix string) T {
	var result T
	if suffix == nilSuffix {
		// Since we decided that 'nil' should equal the empty TypeID, we have to return
		// the empty TypeID when the provided suffix is the the nilSuffix. Otherwise
		// equality will break.
		return result
	}
	PT(&result).init(prefix, suffix)
	return result
}

func subtypePrefix[T Subtype]() string {
	var subtype T
	return subtype.AllowedPrefix()
}
