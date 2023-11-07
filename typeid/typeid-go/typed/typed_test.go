package typed

import (
	"testing"
)

type userID string

const userPrefix = userID("user")

func TestT(t *testing.T) {
	// Arrange
	prefix := userID("user")

	// Act
	builder := T(prefix)

	// Assert
	if builder.prefix != string(prefix) {
		t.Errorf("Expected prefix to be %v, but got %v", prefix, builder.prefix)
	}
}

// TestNew tests the New method that returns a new TypeID with a random suffix.
func TestNew(t *testing.T) {
	builder := T(userID("user"))

	// Act
	tid, err := builder.New()

	// Assert
	// You'd want to mock `untyped.New` here to ensure the behavior
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if tid == nil {
		t.Errorf("Expected a non-nil TypeID, but got nil")
	}

	if tid != nil && tid.Type() != string(userPrefix) {
		t.Errorf("Expected type to be %v, but got %v", userPrefix, tid.Type())
	}
}

var validUserIDs = []string{
	"user_01heknff3res1v8g3q0qyppv08",
	"user_01heknffp5er8tvyecnzkj3sfj",
	"user_01heknfg0ke9gt4mvdpm6gwgc6",
	"user_01heknfgaxf1faykgx5e2hz9xn",
}

var invalidUserIDs = []string{
	"user_1234",
	"user-01heknfgaxf1faykgx5e2hz9xn",
	"account_01heknfgaxf1faykgx5e2hz9xn",
	"",
	"user-",
}

func TestFromString(t *testing.T) {
	builder := T(userID("user"))

	// Test valid IDs
	for _, vID := range validUserIDs {
		tid, err := builder.FromString(vID)
		if err != nil {
			t.Errorf("FromString(%q) returned unexpected error: %v", vID, err)
		}
		if tid == nil {
			t.Errorf("FromString(%q) returned nil, want non-nil TypeID", vID)
		}
		if tid.String() != vID {
			t.Errorf("FromString(%q) = %v, want %v", vID, tid.String(), vID)
		}
		if tid.Type() != string(userPrefix) {
			t.Errorf("FromString(%q) = %v, want %v", vID, tid.Type(), userPrefix)
		}
	}

	// Test invalid IDs
	for _, invID := range invalidUserIDs {
		tid, err := builder.FromString(invID)
		if err == nil {
			t.Errorf("FromString(%q) should have failed but didn't", invID)
		}
		if tid != nil {
			t.Errorf("FromString(%q) returned non-nil TypeID, want nil", invID)
		}
		if tid.String() != "00000000000000000000000000" {
			t.Errorf("FromString(%q) = %v, want %v", invID, tid.String(), "00000000000000000000000000")
		}
		if tid.Type() != "" {
			t.Errorf("FromString(%q) = %v, want %v", invID, tid.Type(), "")
		}
	}
}
