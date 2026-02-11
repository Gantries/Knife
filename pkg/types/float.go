package types

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/gantries/knife/pkg/errors"
)

type Float float64

func (y *Float) Value() float64 {
	return float64(*y)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (y *Float) UnmarshalGQL(v interface{}) error {
	if str, ok := v.(string); ok {
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*y = Float(v)
		return nil
	}
	if jsonV, ok := v.(json.Number); ok {
		v, err := jsonV.Float64()
		if err != nil {
			return err
		}
		*y = Float(v)
		return nil
	}

	switch reflect.ValueOf(v).Kind() {
	case reflect.Int:
		*y = Float(v.(int))
		return nil
	case reflect.Int16:
		*y = Float(v.(int16))
		return nil
	case reflect.Int32:
		*y = Float(v.(int32))
		return nil
	case reflect.Int64:
		*y = Float(v.(int64))
		return nil
	case reflect.Uint:
		*y = Float(v.(uint))
		return nil
	case reflect.Uint8:
		*y = Float(v.(uint8))
		return nil
	case reflect.Uint16:
		*y = Float(v.(uint16))
		return nil
	case reflect.Uint32:
		*y = Float(v.(uint32))
		return nil
	case reflect.Uint64:
		*y = Float(v.(uint64))
		return nil
	case reflect.Float32:
		*y = Float(v.(float32))
		return nil
	case reflect.Float64:
		*y = Float(v.(float64))
		return nil
	default:
	}

	return errors.UnexpectedValueError.E(logger, "type", "float64,string", "value", v)
}

// MarshalGQL implements the graphql.Marshaler interface
func (y Float) MarshalGQL(w io.Writer) {
	_, err := w.Write([]byte(fmt.Sprintf("%f", float64(y))))
	if err != nil {
		return
	}
}

func FloatToString(id Float) string {
	return strconv.FormatFloat(float64(id), 'f', 14, 64)
}
