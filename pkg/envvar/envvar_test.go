package envvar

import "testing"

func TestGet(t *testing.T) {
	const key = "ENVVAR_TEST_GET"

	t.Run("Unset", func(t *testing.T) {
		t.Setenv(key, "")
		if got := Get(key, "default"); got != "default" {
			t.Errorf("Get(%q, %q) = %q, want %q", key, "default", got, "default")
		}
	})

	t.Run("Set", func(t *testing.T) {
		t.Setenv(key, "value")
		if got := Get(key, "default"); got != "value" {
			t.Errorf("Get(%q, %q) = %q, want %q", key, "default", got, "value")
		}
	})

	t.Run("EmptyDefault", func(t *testing.T) {
		t.Setenv(key, "")
		if got := Get(key, ""); got != "" {
			t.Errorf("Get(%q, %q) = %q, want empty", key, "", got)
		}
	})
}

func TestBool(t *testing.T) {
	const key = "ENVVAR_TEST_BOOL"

	tests := []struct {
		val  string
		want bool
	}{
		{"", false},
		{"1", true},
		{"t", true},
		{"T", true},
		{"true", true},
		{"TRUE", true},
		{"True", true},
		{"0", false},
		{"f", false},
		{"false", false},
		{"FALSE", false},
		{"yes", false}, // not understood by strconv.ParseBool
		{"on", false},
		{"garbage", false},
		{" true", false},
	}
	for _, tt := range tests {
		t.Run(tt.val, func(t *testing.T) {
			t.Setenv(key, tt.val)
			if got := Bool(key); got != tt.want {
				t.Errorf("Bool(%q) with value %q = %v, want %v", key, tt.val, got, tt.want)
			}
		})
	}
}
