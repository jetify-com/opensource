package typeid

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

const minTupleStringLength = 40 // (x,00000000-0000-0000-0000-000000000000)

// Scan implements the sql.Scanner interface so the TypeIDs can be read from
// databases transparently. Currently database types that map to string are
// supported.
func (tid *TypeID[P]) Scan(src any) error {
	switch obj := src.(type) {
	case nil:
		return nil
	case string:
		if src == "" {
			return nil
		}
		return tid.UnmarshalText([]byte(obj))
	case []byte:
		if len(obj) == 0 {
			return nil
		}

		// typeid-sql can store TypeIDs as tuples of the form (prefix,uuid).
		if len(obj) < minTupleStringLength || obj[0] != '(' || obj[len(obj)-1] != ')' {
			// TODO: add support for []byte
			// we don't just want to store the full string as a byte array. Instead
			// we should encode using the UUID bytes. We could add support for
			// Binary Marshalling and Unmarshalling at the same time.
			return fmt.Errorf("unsupported format for scan type %T", obj)
		}

		obj = obj[1 : len(obj)-1]
		parts := strings.Split(string(obj), ",")
		if len(parts) != 2 {
			return fmt.Errorf("invalid TypeID format: %s", obj)
		}

		parsedID, err := fromUUID[TypeID[P]](parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("invalid UUID: %s: %w", parts[1], err)
		}

		*tid = parsedID
		return nil
	default:
		return fmt.Errorf("unsupported scan type %T", obj)
	}
}

// Value implements the sql.Valuer interface so that TypeIDs can be written
// to databases transparently. Currently, TypeIDs map to strings.
func (tid TypeID[P]) Value() (driver.Value, error) {
	return tid.String(), nil
}
