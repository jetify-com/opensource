package record

import "log/slog"

// Utilities for converting to slog types.

func (r Record) ToSlog() slog.Record {
	result := slog.NewRecord(r.Time, slog.LevelError, r.Message, r.PC)

	r.Attrs(func(a Attr) bool {
		result.AddAttrs(a.ToSlog())
		return true
	})
	return result
}

func (a Attr) ToSlog() slog.Attr {
	return slog.Attr{Key: a.Key, Value: a.Value.ToSlog()}
}

func (v Value) ToSlog() slog.Value {
	switch v.Kind() {
	case KindString:
		return slog.StringValue(v.String())
	case KindInt64:
		return slog.Int64Value(v.Int64())
	case KindInt:
		return slog.IntValue(v.Int())
	case KindUint64:
		return slog.Uint64Value(v.Uint64())
	case KindDuration:
		return slog.DurationValue(v.Duration())
	case KindTime:
		return slog.TimeValue(v.Time())
	case KindBool:
		return slog.BoolValue(v.Bool())
	case KindFloat64:
		return slog.Float64Value(v.Float64())
	case KindGroup:
		return slog.GroupValue(AttrsToSlog(v.Group())...)
	}
	return slog.AnyValue(v.Any())
}

func AttrsToSlog(attrs []Attr) []slog.Attr {
	slogAttrs := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		slogAttrs[i] = a.ToSlog()
	}
	return slogAttrs
}
