// Package lang contains all language level utilities.
package lang

import (
	"encoding/json"
)

// ToLong converts any convertable value into int64
func ToLong[T any](v interface{}, cvt func(*int64) *T, def *T) *T {
	if v == nil {
		return def
	}
	var n int64
	switch v := v.(type) {
	case string:
		n, _ = json.Number(v).Int64()
	case json.Number:
		n, _ = v.Int64()
	case int:
		n = int64(v)
	case int8:
		n = int64(v)
	case int16:
		n = int64(v)
	case int32:
		n = int64(v)
	case int64:
		n = int64(v)
	case uint:
		n = int64(v) // #nosec G115 - type conversion is safe for typical use cases
	case uint8:
		n = int64(v)
	case uint16:
		n = int64(v)
	case uint32:
		n = int64(v)
	case uint64:
		n = int64(v) // #nosec G115 - type conversion is safe for typical use cases
	case bool:
		n = Ternary(v, int64(1), int64(0))
	case []byte:
		// treat as network order
		buf := v
		m := len(buf) - 1
		n = int64(0)
		for i, b := range buf {
			n |= (int64(b) << (8 * (m - i)))
		}
	default:
		return def
	}
	return cvt(&n)
}
