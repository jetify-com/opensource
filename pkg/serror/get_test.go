package serror

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type custom struct{ value string }
	customValue := custom{"test"}

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	err := New("validation failed",
		"code", 400,
		"time", fixedTime,
		"enabled", true,
		"custom", customValue,
		Group("user",
			Int("id", 123),
			String("name", "alice"),
			Group("address",
				String("city", "Portland"),
				String("state", "OR"),
			),
		),
		Group("request",
			String("path", "/api/users"),
			Int64("size", 3000000000),
		),
	)

	tests := []struct {
		path     string
		want     any
		wantKind Kind
	}{
		{"code", 400, KindInt},
		{"time", fixedTime, KindTime},
		{"enabled", true, KindBool},
		{"user.id", 123, KindInt},
		{"user.name", "alice", KindString},
		{"user.address.city", "Portland", KindString},
		{"user.address.state", "OR", KindString},
		{"request.path", "/api/users", KindString},
		{"request.size", int64(3000000000), KindInt64},
		{"custom", customValue, KindAny},
		{"nonexistent", nil, KindAny},
		{"user.nonexistent", nil, KindAny},
		{"user.address.nonexistent", nil, KindAny},
		// Test getting entire groups
		{"user", []Attr{
			Int("id", 123),
			String("name", "alice"),
			Group("address",
				String("city", "Portland"),
				String("state", "OR"),
			),
		}, KindGroup},
		{"user.address", []Attr{
			String("city", "Portland"),
			String("state", "OR"),
		}, KindGroup},
		{"request", []Attr{
			String("path", "/api/users"),
			Int64("size", 3000000000),
		}, KindGroup},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			got := err.Get(test.path)
			if test.want == nil {
				assert.Empty(t, got, "expected empty value for %q", test.path)
				return
			}

			assert.Equal(t, test.wantKind, got.Kind(), "unexpected kind for %q", test.path)
			if test.wantKind == KindGroup {
				assert.ElementsMatch(t, test.want.([]Attr), got.Group(),
					"unexpected group attributes for %q", test.path)
			} else {
				assert.Equal(t, test.want, got.Any(), "unexpected value for %q", test.path)
			}
		})
	}
}
