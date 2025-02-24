package serror

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return fixedTime }
	defer func() { timeNow = time.Now }()

	tests := []struct {
		name     string
		err      Error
		wantJSON string
	}{
		{
			name: "bool field",
			err:  New("test bool", Bool("field", true)),
			wantJSON: `{
				"message": "test bool",
				"field": true
			}`,
		},
		{
			name: "duration field",
			err:  New("test duration", Duration("field", 5*time.Second)),
			wantJSON: `{
				"message": "test duration",
				"field": {
					"kind": "Duration",
					"value": "5s"
				}
			}`,
		},
		{
			name: "float64 field",
			err:  New("test float", Float64("field", 123.456)),
			wantJSON: `{
				"message": "test float",
				"field": 123.456
			}`,
		},
		{
			name: "int64 field",
			err:  New("test int", Int("field", -123)),
			wantJSON: `{
				"message": "test int",
				"field": {
					"kind": "Int64",
					"value": -123
				}
			}`,
		},
		{
			name: "string field",
			err:  New("test string", String("field", "hello")),
			wantJSON: `{
				"message": "test string",
				"field": "hello"
			}`,
		},
		{
			name: "time field",
			err:  New("test time", Time("field", fixedTime)),
			wantJSON: `{
				"message": "test time",
				"field": {
					"kind": "Time",
					"value": "2024-01-01T00:00:00Z"
				}
			}`,
		},
		{
			name: "uint64 field",
			err:  New("test uint", Uint64("field", 123)),
			wantJSON: `{
				"message": "test uint",
				"field": {
					"kind": "Uint64",
					"value": 123
				}
			}`,
		},
		{
			name: "group field",
			err: New("test group",
				Group("outer",
					String("inner", "value"),
					Int("number", 42),
				),
			),
			wantJSON: `{
				"message": "test group",
				"outer": {
					"inner": "value",
					"number": {
						"kind": "Int64",
						"value": 42
					}
				}
			}`,
		},
		{
			name: "error with cause",
			err: Wrap(
				New("inner error", String("detail", "info")),
				"outer error",
				Int("code", 500),
			),
			wantJSON: `{
				"message": "outer error",
				"code": {
					"kind": "Int64",
					"value": 500
				},
				"cause": {
					"message": "inner error",
					"detail": "info"
				}
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.err)
			require.NoError(t, err)

			assert.JSONEq(t, test.wantJSON, string(data), "json mismatch")

			var gotErr Error
			require.NoError(t, json.Unmarshal(data, &gotErr))

			wantMap := test.err.toMap()
			gotMap := gotErr.toMap()
			assert.Equal(t, wantMap, gotMap, "deserialized structure mismatch")
		})
	}
}
