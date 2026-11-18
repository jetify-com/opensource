package typeid

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// For nullable TypeID columns, use sql.Null[TypeID].

const minTupleStringLength = 40 // (x,00000000-0000-0000-0000-000000000000)

// Scan implements the sql.Scanner interface so the TypeIDs can be read from
// databases transparently. Currently database types that map to string are
// supported.
func (tid *TypeID) Scan(src any) error {
	switch obj := src.(type) {
	case nil:
		return &validationError{
			Message: "cannot scan NULL into TypeID",
		}
	case string:
		if obj == "" {
			return &validationError{
				Message: "cannot scan empty string into TypeID",
			}
		}
		return tid.UnmarshalText([]byte(obj))
	case []byte:
		// typeid-sql can store TypeIDs as tuples of the form (prefix,uuid).
		if len(obj) < minTupleStringLength || obj[0] != '(' || obj[len(obj)-1] != ')' {
			return &validationError{
				Message: fmt.Sprintf("unsupported scan type %T", obj),
			}
		}

		parts := strings.Split(string(obj[1:len(obj)-1]), ",")
		if len(parts) != 2 {
			return &validationError{
				Message: fmt.Sprintf("invalid TypeID tuple %q", obj),
			}
		}

		parsed, err := FromUUID(parts[0], parts[1])
		if err != nil {
			return err
		}

		*tid = parsed
		return nil
	default:
		return &validationError{
			Message: fmt.Sprintf("unsupported scan type %T", obj),
		}
	}
}

// Value implements the sql.Valuer interface so that TypeIDs can be written
// to databases transparently. Currently, TypeIDs map to strings.
func (tid TypeID) Value() (driver.Value, error) {
	return tid.String(), nil
}
