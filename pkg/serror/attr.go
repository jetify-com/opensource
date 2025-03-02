package serror

import (
	"time"

	"go.jetify.com/pkg/serror/internal/record"
)

const badKey = "!BADKEY"

// Kind is the kind of a [Value].
type Kind = record.Kind

const (
	KindAny       = record.KindAny
	KindBool      = record.KindBool
	KindDuration  = record.KindDuration
	KindFloat64   = record.KindFloat64
	KindInt64     = record.KindInt64
	KindString    = record.KindString
	KindTime      = record.KindTime
	KindUint64    = record.KindUint64
	KindGroup     = record.KindGroup
	KindLogValuer = record.KindLogValuer
)

// An Attr is a key-value pair.
type Attr = record.Attr

// String returns an Attr for a string value.
func String(key, value string) Attr {
	return record.String(key, value)
}

// Int64 returns an Attr for an int64.
func Int64(key string, value int64) Attr {
	return record.Int64(key, value)
}

// Int converts an int to an int64 and returns
// an Attr with that value.
func Int(key string, value int) Attr {
	return record.Int(key, value)
}

// Uint64 returns an Attr for a uint64.
func Uint64(key string, v uint64) Attr {
	return record.Uint64(key, v)
}

// Float64 returns an Attr for a floating-point number.
func Float64(key string, v float64) Attr {
	return record.Float64(key, v)
}

// Bool returns an Attr for a bool.
func Bool(key string, v bool) Attr {
	return record.Bool(key, v)
}

// Time returns an Attr for a [time.Time].
// It discards the monotonic portion.
func Time(key string, v time.Time) Attr {
	return record.Time(key, v)
}

// Duration returns an Attr for a [time.Duration].
func Duration(key string, v time.Duration) Attr {
	return record.Duration(key, v)
}

// Group returns an Attr for a Group [Value].
// The first argument is the key; the remaining arguments
// are converted to Attrs as in [Logger.Log].
//
// Use Group to collect several key-value pairs under a single
// key on a log line, or as the result of LogValue
// in order to log a single value as multiple Attrs.
func Group(key string, args ...any) Attr {
	return record.Group(key, args...)
}

// Any returns an Attr for the supplied value.
// See [AnyValue] for how values are treated.
func Any(key string, value any) Attr {
	return record.Any(key, value)
}
