package serror

import (
	"log/slog"
	"strings"
)

// Get returns the Value for a given path using dot notation.
// For example, "code" returns the Value under the key "code",
// while "user.id" returns the Value under key "id" in group "user".
// If the path doesn't exist, it returns a nil Value.
func (e Error) Get(path string) Value {
	parts := strings.Split(path, ".")

	result := slog.Value{}
	e.record.Attrs(func(attr slog.Attr) bool {
		if value := findValue(attr, parts); !value.Equal(slog.Value{}) {
			result = value
			return false // stop iteration
		}
		return true
	})

	return result
}

// findValue recursively searches for a value in an Attr using the path parts
func findValue(attr slog.Attr, parts []string) slog.Value {
	if len(parts) == 0 {
		return slog.Value{}
	}

	// Check if this attribute matches the first part
	if attr.Key != parts[0] {
		return slog.Value{}
	}

	// If this is the last part, return the value
	if len(parts) == 1 {
		return attr.Value
	}

	// If we have more parts, this must be a group
	if attr.Value.Kind() != slog.KindGroup {
		return slog.Value{}
	}

	// Search through the group's attributes
	for _, groupAttr := range attr.Value.Group() {
		if value := findValue(groupAttr, parts[1:]); !value.Equal(slog.Value{}) {
			return value
		}
	}

	return slog.Value{}
}
