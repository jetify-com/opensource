package serror

import (
	"bytes"
	"errors"
	"log/slog"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstructors(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return fixedTime }
	defer func() { timeNow = time.Now }()

	baseErr := New("base error")

	tests := []struct {
		name string
		args []any
		expectedError
	}{
		{
			name:          "basic message without attrs",
			expectedError: expectedError{message: "test error"},
		},
		{
			name: "basic types",
			args: []any{"str", "value", "int", 42, "bool", true, "float", 3.14},
			expectedError: expectedError{
				message:  "test error",
				wantKeys: []string{"str", "int", "bool", "float"},
				wantVals: []string{"value", "42", "true", "3.14"},
			},
		},
		{
			name: "time values",
			args: []any{"time", fixedTime, "duration", time.Hour},
			expectedError: expectedError{
				message:  "test error",
				wantKeys: []string{"time", "duration"},
				wantVals: []string{fixedTime.String(), "1h0m0s"},
			},
		},
		{
			name: "direct attrs",
			args: []any{Int("count", 5), String("msg", "hello"), Bool("ok", true)},
			expectedError: expectedError{
				message:  "test error",
				wantKeys: []string{"count", "msg", "ok"},
				wantVals: []string{"5", "hello", "true"},
			},
		},
		{
			name: "groups",
			args: []any{
				Group("user",
					Int("id", 123),
					String("name", "alice"),
				),
				Group("metrics",
					Int("count", 100),
					Duration("latency", time.Second),
				),
			},
			expectedError: expectedError{
				message:  "test error",
				wantKeys: []string{"user.id", "user.name", "metrics.count", "metrics.latency"},
				wantVals: []string{"123", "alice", "100", "1s"},
			},
		},
		{
			name: "mixed types",
			args: []any{"direct", "value", Int("count", 42), Group("group", "name", "test", "ok", true)},
			expectedError: expectedError{
				message:  "test error",
				wantKeys: []string{"direct", "count", "group.name", "group.ok"},
				wantVals: []string{"value", "42", "test", "true"},
			},
		},
		{
			name: "invalid args",
			args: []any{123, 44},
			expectedError: expectedError{
				message:  "test error",
				wantKeys: []string{badKey, badKey},
				wantVals: []string{"123", "44"},
			},
		},
	}

	t.Run("New", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				err := New(test.message, test.args...)
				assertError(t, err, test.expectedError)
			})
		}
	})

	t.Run("Wrap", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				want := test.expectedError
				want.cause = baseErr
				err := Wrap(baseErr, test.message, test.args...)
				assertError(t, err, want)
				assert.Equal(t, baseErr, errors.Unwrap(err))
			})
		}
	})
}

func TestCallLocation(t *testing.T) {
	// Create helper function to get a clean line number for comparison
	newWithLine := func() (Error, int) {
		_, _, line, _ := runtime.Caller(0)
		return New("test error"), line + 1
	}
	wrapWithLine := func(err error) (Error, int) {
		_, _, line, _ := runtime.Caller(0)
		return Wrap(err, "wrapped"), line + 1
	}

	t.Run("New captures caller's location", func(t *testing.T) {
		err, line := newWithLine()
		frame, _ := runtime.CallersFrames([]uintptr{err.record.PC}).Next()
		assert.Contains(t, frame.File, "serror_test.go")
		assert.Equal(t, line, frame.Line)
	})

	t.Run("Wrap captures caller's location", func(t *testing.T) {
		base := New("base error")
		err, line := wrapWithLine(base)
		frame, _ := runtime.CallersFrames([]uintptr{err.record.PC}).Next()
		assert.Contains(t, frame.File, "serror_test.go")
		assert.Equal(t, line, frame.Line)
	})
}

func TestErrorString(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{
			name: "basic message",
			err:  New("failed"),
			want: "failed",
		},
		{
			name: "with attributes",
			err:  New("failed", "user", "alice", "count", 42),
			want: "failed user=alice count=42",
		},
		{
			name: "with cause",
			err:  Wrap(New("root cause"), "operation failed", "id", 123),
			want: "operation failed id=123 cause=\"root cause\"",
		},
		{
			name: "with nested cause and attrs",
			err: Wrap(
				Wrap(New("root", "file", "config.json"),
					"load failed",
					"type", "config"),
				"init failed",
				"service", "api"),
			want: "init failed service=api cause=\"load failed type=config cause=\\\"root file=config.json\\\"\"",
		},
		{
			name: "with group",
			err: New("failed",
				Group("user",
					Int("id", 123),
					String("name", "alice"),
				),
			),
			want: "failed user.id=123 user.name=alice",
		},
		{
			name: "with newlines in message",
			err:  New("line1\nline2"),
			want: "line1\nline2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.err.Error()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestWith(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return fixedTime }
	defer func() { timeNow = time.Now }()

	tests := []struct {
		name     string
		baseErr  Error
		withArgs []any
		want     string
	}{
		{
			name:     "add simple attributes",
			baseErr:  New("base error", "initial", "value"),
			withArgs: []any{"added", "attr", "num", 42},
			want:     "base error initial=value added=attr num=42",
		},
		{
			name:    "add group",
			baseErr: New("base error", "user", "alice"),
			withArgs: []any{Group("metrics",
				Int("count", 100),
				String("latency", "50ms"),
			)},
			want: "base error user=alice metrics.count=100 metrics.latency=50ms",
		},
		{
			name: "add to error with cause",
			baseErr: Wrap(
				New("root cause"),
				"operation failed",
				"initial", "value",
			),
			withArgs: []any{"status", 500},
			want:     "operation failed initial=value status=500 cause=\"root cause\"",
		},
		{
			name:     "empty with args",
			baseErr:  New("base error", "key", "value"),
			withArgs: []any{},
			want:     "base error key=value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Keep original for comparison
			original := test.baseErr.Error()

			// Create new error with additional attributes
			withErr := test.baseErr.With(test.withArgs...)

			// Verify original is unchanged
			assert.Equal(t, original, test.baseErr.Error(), "original error should not be modified")

			// Verify new error has expected string representation
			assert.Equal(t, test.want, withErr.Error())

			// Verify cause is preserved
			assert.Equal(t, test.baseErr.cause, withErr.cause)
		})
	}
}

func TestLogValue(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		// Remove time to make output predictable
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	tests := []struct {
		name     string
		err      Error
		wantLogs string
	}{
		{
			name: "basic error",
			err:  New("operation failed", "status", 500),
			wantLogs: `level=ERROR msg="log error" error.status=500 ` +
				`error.msg="operation failed"` + "\n",
		},
		{
			name: "error with cause",
			err: Wrap(
				New("database error", "table", "users"),
				"query failed",
				"query_id", "abc123",
			),
			wantLogs: `level=ERROR msg="log error" error.query_id=abc123 ` +
				`error.msg="query failed" error.cause.table=users ` +
				`error.cause.msg="database error"` + "\n",
		},
		{
			name: "error with groups",
			err: New("validation failed",
				Group("user",
					String("id", "123"),
					String("name", "alice"),
				),
				Group("request",
					Int("status", 400),
					String("path", "/api/v1/users"),
				),
			),
			wantLogs: `level=ERROR msg="log error" error.user.id=123 ` +
				`error.user.name=alice error.request.status=400 ` +
				`error.request.path=/api/v1/users error.msg="validation failed"` + "\n",
		},
		{
			name: "deeply nested error",
			err: Wrap(
				Wrap(
					New("root cause", "file", "config.json"),
					"load failed",
					"type", "config",
				),
				"init failed",
				"service", "api",
			),
			wantLogs: `level=ERROR msg="log error" error.service=api ` +
				`error.msg="init failed" error.cause.type=config ` +
				`error.cause.msg="load failed" error.cause.cause.file=config.json ` +
				`error.cause.cause.msg="root cause"` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf.Reset()
			logger.Error("log error", "error", test.err)
			assert.Equal(t, test.wantLogs, buf.String())
		})
	}
}

type expectedError struct {
	message  string
	cause    error
	wantKeys []string
	wantVals []string
}

func assertError(t *testing.T, err Error, want expectedError) {
	assert.Equal(t, want.message, err.record.Message)
	assert.Equal(t, want.cause, err.cause)

	if len(want.wantKeys) == 0 {
		return
	}

	var flattened []Attr
	err.record.Attrs(func(a slog.Attr) bool {
		flattened = append(flattened, flattenAttrs(a)...)
		return true
	})

	assert.Len(t, flattened, len(want.wantKeys), "unexpected number of attributes")
	for i, key := range want.wantKeys {
		assert.Equal(t, key, flattened[i].Key, "key mismatch at index %d", i)
		assert.Equal(t, want.wantVals[i], flattened[i].Value.String(), "value mismatch at index %d", i)
	}
}

func flattenAttrs(attr slog.Attr) []Attr {
	if attr.Value.Kind() == slog.KindGroup {
		var results []Attr
		for _, ga := range attr.Value.Group() {
			for _, sub := range flattenAttrs(ga) {
				if attr.Key != "" {
					sub.Key = attr.Key + "." + sub.Key
				}
				results = append(results, sub)
			}
		}
		return results
	}
	return []Attr{{Key: attr.Key, Value: attr.Value}}
}
