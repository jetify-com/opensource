package optional

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyIndependence(t *testing.T) {
	tests := []struct {
		name string
		init Option[int]
		mod  string // JSON to modify copy
	}{
		{
			name: "some to some",
			init: Some(42),
			mod:  "99",
		},
		{
			name: "some to none",
			init: Some(42),
			mod:  "null",
		},
		{
			name: "none to some",
			init: None[int](),
			mod:  "99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create original option
			original := tt.init

			// Make a copy
			copy := original

			// Save original's state
			origVal, origPresent := original.Get()

			// Modify copy via UnmarshalJSON
			err := json.Unmarshal([]byte(tt.mod), &copy)
			require.NoError(t, err)

			// Verify original is unchanged
			newVal, newPresent := original.Get()
			assert.Equal(t, origVal, newVal)
			assert.Equal(t, origPresent, newPresent)

			// Verify copy was changed
			if tt.mod == "null" {
				assert.True(t, copy.IsNone())
			} else {
				val, ok := copy.Get()
				assert.True(t, ok)
				assert.Equal(t, 99, val)
			}
		})
	}
}

func TestReferenceTypeContainment(t *testing.T) {
	// Create a map as the inner value
	m := map[string]int{"key": 1}

	// Create option and copy
	original := Some(m)
	copy := original

	// Get maps from both options
	origMap, origOk := original.Get()
	require.True(t, origOk)

	copyMap, copyOk := copy.Get()
	require.True(t, copyOk)

	// Verify we have the same map instance
	assert.Equal(t, origMap, copyMap)

	// Modify map through copy
	copyMap["key"] = 2

	// Verify changes visible through original (reference semantics of map)
	checkMap, _ := original.Get()
	assert.Equal(t, 2, checkMap["key"])

	// Modify the copy's Option wrapper
	err := json.Unmarshal([]byte(`null`), &copy)
	require.NoError(t, err)

	// Original should be unaffected
	_, stillPresent := original.Get()
	assert.True(t, stillPresent)
	assert.False(t, original.IsNone())

	// Copy should be None
	assert.True(t, copy.IsNone())
}

func TestMethodImmutability(t *testing.T) {
	tests := []struct {
		name   string
		option Option[int]
	}{
		{
			name:   "some value",
			option: Some(42),
		},
		{
			name:   "none value",
			option: None[int](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save initial state
			initialValue, initialPresent := tt.option.Get()
			initialNone := tt.option.IsNone()

			// Call every non-mutating method
			_ = tt.option.GetOrElse(0)
			_ = tt.option.IsNone()
			_ = tt.option.String()
			_ = tt.option.GoString()
			_, _ = tt.option.MarshalJSON()
			_ = tt.option.IsZero()

			// Verify none of the methods changed the option
			finalValue, finalPresent := tt.option.Get()
			finalNone := tt.option.IsNone()

			assert.Equal(t, initialValue, finalValue)
			assert.Equal(t, initialPresent, finalPresent)
			assert.Equal(t, initialNone, finalNone)
		})
	}
}
