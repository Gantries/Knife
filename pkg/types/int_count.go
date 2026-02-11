package types

import (
	"io"
	"reflect"
	"strconv"

	"github.com/gantries/knife/pkg/errors"
)

type IntCount int64

func (y *IntCount) UnmarshalGQL(v interface{}) error {
	if s, ok := v.(string); ok {
		if i64, err := strconv.ParseInt(s, 10, 64); err != nil {
			return err
		} else {
			*y = IntCount(i64)
			return nil
		}
	}

	switch reflect.ValueOf(v).Kind() {
	case reflect.Int:
		*y = IntCount(v.(int))
		return nil
	case reflect.Int16:
		*y = IntCount(v.(int16))
		return nil
	case reflect.Int32:
		*y = IntCount(v.(int32))
		return nil
	case reflect.Int64:
		*y = IntCount(v.(int64))
		return nil
	case reflect.Uint:
		*y = IntCount(v.(uint)) // #nosec G115 - type conversion is safe for typical use cases
		return nil
	case reflect.Uint8:
		*y = IntCount(v.(uint8))
		return nil
	case reflect.Uint16:
		*y = IntCount(v.(uint16))
		return nil
	case reflect.Uint32:
		*y = IntCount(v.(uint32))
		return nil
	case reflect.Uint64:
		*y = IntCount(v.(uint64)) // #nosec G115 - type conversion is safe for typical use cases
		return nil
	default:
	}

	return errors.UnexpectedValueError.E(logger, "target", "int-count", "value", v)
}

func (y IntCount) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.FormatInt(int64(y), 10)))
}
