package errors

import (
	stderrors "errors"
	"fmt"
	"io/fs"
	"strings"
	"testing"
)

var errSentinel = New("sentinel")

func TestNewAndErrorf(t *testing.T) {
	if got := New("boom").Error(); got != "boom" {
		t.Errorf("New(%q).Error() = %q", "boom", got)
	}
	if got := Errorf("value %d is %s", 42, "bad").Error(); got != "value 42 is bad" {
		t.Errorf("Errorf().Error() = %q", got)
	}
}

func TestNewHasStackTrace(t *testing.T) {
	// pkg/errors errors print a stack trace with %+v.
	got := fmt.Sprintf("%+v", New("boom"))
	if !strings.Contains(got, "TestNewHasStackTrace") {
		t.Errorf("expected stack trace to mention the test function, got:\n%s", got)
	}
}

func TestNilPassthrough(t *testing.T) {
	tests := map[string]error{
		"Wrap":         Wrap(nil, "msg"),
		"Wrapf":        Wrapf(nil, "msg %d", 1),
		"WithMessage":  WithMessage(nil, "msg"),
		"WithMessagef": WithMessagef(nil, "msg %d", 1),
		"WithStack":    WithStack(nil),
		"Cause":        Cause(nil),
		"Unwrap":       Unwrap(nil),
		"Join":         Join(),
		"JoinNils":     Join(nil, nil),
	}
	for name, err := range tests {
		if err != nil {
			t.Errorf("%s(nil) = %v, want nil", name, err)
		}
	}
}

func TestWrapping(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantMsg string
	}{
		{"Wrap", Wrap(errSentinel, "ctx"), "ctx: sentinel"},
		{"Wrapf", Wrapf(errSentinel, "ctx %d", 7), "ctx 7: sentinel"},
		{"WithMessage", WithMessage(errSentinel, "ctx"), "ctx: sentinel"},
		{"WithMessagef", WithMessagef(errSentinel, "ctx %s", "x"), "ctx x: sentinel"},
		{"WithStack", WithStack(errSentinel), "sentinel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("Error() = %q, want %q", got, tt.wantMsg)
			}
			if !Is(tt.err, errSentinel) {
				t.Error("Is(err, sentinel) = false, want true")
			}
			if got := Cause(tt.err); got != errSentinel { //nolint:errorlint // Cause must return the exact root error
				t.Errorf("Cause(err) = %v, want sentinel", got)
			}
			if Unwrap(tt.err) == nil {
				t.Error("Unwrap(err) = nil, want non-nil")
			}
		})
	}
}

func TestCauseOfUnwrappedError(t *testing.T) {
	plain := stderrors.New("plain")
	if got := Cause(plain); got != plain { //nolint:errorlint // Cause must return the exact root error
		t.Errorf("Cause(plain) = %v, want plain", got)
	}
}

func TestIsWithStdlibWrapping(t *testing.T) {
	err := fmt.Errorf("outer: %w", Wrap(fs.ErrNotExist, "inner"))
	if !Is(err, fs.ErrNotExist) {
		t.Error("Is() should see through mixed stdlib and pkg/errors wrapping")
	}
	if Is(err, fs.ErrExist) {
		t.Error("Is() matched an unrelated error")
	}
}

type customErr struct{ code int }

func (e *customErr) Error() string { return fmt.Sprintf("code %d", e.code) }

func TestAs(t *testing.T) {
	err := Wrap(&customErr{code: 3}, "ctx")

	var target *customErr
	if !As(err, &target) {
		t.Fatal("As() = false, want true")
	}
	if target.code != 3 {
		t.Errorf("target.code = %d, want 3", target.code)
	}

	var pathErr *fs.PathError
	if As(err, &pathErr) {
		t.Error("As() matched an unrelated error type")
	}
}

func TestJoin(t *testing.T) {
	a, b := New("a"), New("b")
	err := Join(a, nil, b)
	if err == nil {
		t.Fatal("Join() = nil, want non-nil")
	}
	if got := err.Error(); got != "a\nb" {
		t.Errorf("Join().Error() = %q, want %q", got, "a\nb")
	}
	if !Is(err, a) || !Is(err, b) {
		t.Error("joined error should match each of its members")
	}
}
