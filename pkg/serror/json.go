package serror

import (
	"encoding/json"
	"errors"
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
	e.record.Attrs(func(a Attr) bool {
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
// Special types (Int, Int64, Uint64, Duration, Time) are wrapped to preserve their type information.
func addAttr(attrMap map[string]any, attr Attr) {
	if attr.Value.Kind() == KindGroup {
		groupMap := make(map[string]any)
		for _, groupAttr := range attr.Value.Group() {
			addAttr(groupMap, groupAttr)
		}
		attrMap[attr.Key] = groupMap
		return
	}

	switch attr.Value.Kind() {
	case KindInt, KindInt64, KindUint64:
		attrMap[attr.Key] = typedValue{
			Kind:  attr.Value.Kind().String(),
			Value: attr.Value.Any(),
		}
	case KindDuration:
		attrMap[attr.Key] = typedValue{
			Kind:  attr.Value.Kind().String(),
			Value: attr.Value.Duration().String(),
		}
	case KindTime:
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
	var attrs []any
	for k, v := range jsonMap {
		if k == "message" || k == "cause" {
			continue
		}

		if vMap, ok := v.(map[string]any); ok {
			// Check if it's a typed value
			if kind, hasKind := vMap["kind"].(string); hasKind && vMap["value"] != nil {
				switch kind {
				case "Int":
					if num, ok := toNumber(vMap["value"]); ok {
						attrs = append(attrs, k, int(num))
						continue
					}
				case "Int64":
					if num, ok := toNumber(vMap["value"]); ok {
						attrs = append(attrs, k, int64(num))
						continue
					}
				case "Uint64":
					if num, ok := toNumber(vMap["value"]); ok {
						attrs = append(attrs, k, uint64(num))
						continue
					}
				case "Duration":
					if s, ok := vMap["value"].(string); ok {
						if d, err := time.ParseDuration(s); err == nil {
							attrs = append(attrs, k, d)
							continue
						}
					}
				case "Time":
					if s, ok := vMap["value"].(string); ok {
						if t, err := time.Parse(time.RFC3339, s); err == nil {
							attrs = append(attrs, k, t)
							continue
						}
					}
				}
			}

			// It's a regular map/group - handle it as a group
			groupAttrs := mapToAttrs(vMap)
			attrs = append(attrs, Group(k, groupAttrs...))
		} else {
			// It's a regular attribute
			attrs = append(attrs, k, v)
		}
	}
	return attrs
}

// toNumber converts various JSON number formats to float64
func toNumber(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint64:
		return float64(n), true
	case json.Number:
		if f, err := n.Float64(); err == nil {
			return f, true
		}
	}
	return 0, false
}
