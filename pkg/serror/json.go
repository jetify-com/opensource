package serror

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"
)

// typedValue represents values that need type preservation in JSON.
// This wrapper preserves type information for values that would otherwise
// lose their exact type during JSON encoding/decoding.
type typedValue struct {
	Kind  string `json:"kind"`
	Value any    `json:"value"`
}

// MarshalJSON implements json.Marshaler, allowing Error values to be encoded as JSON.
// The JSON representation includes the error message, cause (if any), and all attributes.
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.toMap())
}

// UnmarshalJSON implements json.Unmarshaler, allowing Error values to be decoded from JSON.
// The JSON representation must include at least an error message. Cause and attributes are optional.
func (e *Error) UnmarshalJSON(data []byte) error {
	var jsonMap map[string]any
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return err
	}
	*e = fromMap(jsonMap)
	return nil
}

// toMap converts an Error to a map representation suitable for JSON marshaling.
// The map includes the error message, cause (if present), and all attributes.
func (e Error) toMap() map[string]any {
	jsonMap := make(map[string]any)
	jsonMap["message"] = e.record.Message

	if e.cause != nil {
		// If the cause is another Error, convert it recursively
		// Otherwise, use the error's string representation
		if causeErr, ok := e.cause.(Error); ok { //nolint:errorlint
			jsonMap["cause"] = causeErr.toMap()
		} else {
			jsonMap["cause"] = e.cause.Error()
		}
	}

	// Walk through all attributes and add them to the map
	e.record.Attrs(func(a slog.Attr) bool {
		addAttr(jsonMap, a)
		return true
	})

	return jsonMap
}

// fromMap creates an Error from a map representation.
// The map must contain a "message" key. "cause" and attribute keys are optional.
func fromMap(jsonMap map[string]any) Error {
	msg := extractMessage(jsonMap)
	cause := extractCause(jsonMap)
	e := new(timeNow(), msg, cause, 0)

	if attrs := mapToAttrs(jsonMap); len(attrs) > 0 {
		e.add(attrs...)
	}

	return e
}

// extractMessage extracts and removes the message from the map.
// Returns an empty string if no message is present.
func extractMessage(jsonMap map[string]any) string {
	msg, _ := jsonMap["message"].(string)
	delete(jsonMap, "message")
	return msg
}

// extractCause extracts and removes the cause from the map.
// Returns nil if no cause is present.
func extractCause(jsonMap map[string]any) error {
	if c, ok := jsonMap["cause"]; ok {
		delete(jsonMap, "cause")
		switch v := c.(type) {
		case map[string]any:
			return fromMap(v)
		case string:
			return errors.New(v)
		}
	}
	return nil
}

// addAttr adds an Attr to a map, handling groups recursively.
// Special types (Int64, Uint64, Duration, Time) are wrapped to preserve their type information.
func addAttr(attrMap map[string]any, attr slog.Attr) {
	if attr.Value.Kind() == slog.KindGroup {
		groupMap := make(map[string]any)
		for _, groupAttr := range attr.Value.Group() {
			addAttr(groupMap, groupAttr)
		}
		attrMap[attr.Key] = groupMap
		return
	}

	switch attr.Value.Kind() {
	case slog.KindInt64, slog.KindUint64:
		attrMap[attr.Key] = typedValue{
			Kind:  attr.Value.Kind().String(),
			Value: attr.Value.Any(),
		}
	case slog.KindDuration:
		attrMap[attr.Key] = typedValue{
			Kind:  attr.Value.Kind().String(),
			Value: attr.Value.Duration().String(),
		}
	case slog.KindTime:
		attrMap[attr.Key] = typedValue{
			Kind:  attr.Value.Kind().String(),
			Value: attr.Value.Time().Format(time.RFC3339),
		}
	default:
		attrMap[attr.Key] = attr.Value.Any()
	}
}

// mapToAttrs converts a map to a slice of attribute arguments.
// It handles nested groups and unwraps specially encoded values.
func mapToAttrs(jsonMap map[string]any) []any {
	attrs := make([]any, 0, len(jsonMap)*2)
	for key, v := range jsonMap {
		if groupMap, ok := v.(map[string]any); ok {
			if val, ok := handleTypedValue(groupMap); ok {
				attrs = append(attrs, key, val)
				continue
			}
			groupAttrs := mapToAttrs(groupMap)
			attrs = append(attrs, Group(key, groupAttrs...))
			continue
		}
		attrs = append(attrs, key, v)
	}
	return attrs
}

// handleTypedValue attempts to convert a map to a typed value.
// Returns the value and true if the map represents a supported type
// (Int64, Uint64, Duration, or Time), otherwise returns nil and false.
func handleTypedValue(groupMap map[string]any) (any, bool) {
	kind, ok := groupMap["kind"].(string)
	if !ok || groupMap["value"] == nil {
		return nil, false
	}

	switch kind {
	case "Int64":
		n, ok := groupMap["value"].(float64)
		if !ok {
			return nil, false
		}
		return int64(n), true
	case "Uint64":
		n, ok := groupMap["value"].(float64)
		if !ok {
			return nil, false
		}
		return uint64(n), true
	case "Duration":
		s, ok := groupMap["value"].(string)
		if !ok {
			return nil, false
		}
		d, err := time.ParseDuration(s)
		if err != nil {
			return nil, false
		}
		return d, true
	case "Time":
		s, ok := groupMap["value"].(string)
		if !ok {
			return nil, false
		}
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, false
		}
		return t, true
	default:
		return nil, false
	}
}
