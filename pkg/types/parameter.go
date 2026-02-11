package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/lang"
)

type ParameterType string

const (
	TypeBool      ParameterType = "bool"
	TypeDouble    ParameterType = "double"
	TypeIdentity  ParameterType = "id"
	TypeInt       ParameterType = "int"
	TypeGroup     ParameterType = "group"
	TypeJson      ParameterType = "json"
	TypeState     ParameterType = "state"
	TypeString    ParameterType = "string"
	TypeText      ParameterType = "text"
	TypeTimestamp ParameterType = "timestamp"
	TypeArray     ParameterType = "array"
)

var (
	once      sync.Once
	location  *time.Location
	layoutStr = "2006-01-02 15:04:05" // 固化的layout
)

func initLocation() {
	var err error
	location, err = time.LoadLocation("Local")
	if err != nil {
		log.Fatalf("Error loading local location: %v", err)
	}
}

func (y ParameterType) String() string {
	return string(y)
}

func (y *ParameterType) UnmarshalGQL(v interface{}) error {
	if s, ok := v.(string); ok {
		*y = ParameterType(s)
		return nil
	}
	return errors.UnexpectedValueError.E(logger, "target", "parameter-type", "value", v)
}

func (y ParameterType) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(fmt.Sprintf(`"%s"`, y.String())))
}
func (y ParameterType) Value() (driver.Value, error) {
	return string(y), nil
}

// Scan 实现了 sql.Scanner 接口
func (y *ParameterType) Scan(value interface{}) error {
	switch value := value.(type) {
	case []uint8:
		*y = ParameterType(value)
		return nil
	case string:
		*y = ParameterType(value)
		return nil
	}
	return errors.UnexpectedValueError.E(logger, "type", "!string/![]uint8", "value", value)
}

func (c ParameterType) Equal(other *ParameterType) bool {
	return other != nil && c == *other
}

func (y ParameterType) To(val interface{}) (v string, err error) {
	if nil == val {
		return
	}
	switch y {
	case TypeBool:
		if b, ok := val.(float64); ok {
			return lang.Ternary(b == 1, "true", "false"), nil
		}
		if b, ok := val.(bool); ok {
			return lang.Ternary(b, "true", "false"), nil
		}
		if b, ok := val.(string); ok {
			return lang.Ternary(b == "1" || b == "true", "true", "false"), nil
		}
		logger.Error("Unable to convert to bool", "val", val)
		return
	case TypeDouble:
		// DOUBLE need to be stored in string to prevent from losing precision
		if f, ok := val.(string); ok {
			return f, nil
		}
		if f, ok := val.(float64); ok {
			return strconv.FormatFloat(f, 'f', -1, 64), nil
		}
		logger.Error("Unable to convert to double(float)", "val", val)
		return
	case TypeIdentity:
		if i, ok := val.(int64); ok {
			return strconv.FormatInt(i, 10), nil
		}
		if i, ok := val.(uint64); ok {
			return strconv.FormatUint(i, 10), nil
		}
		if i, ok := val.(float64); ok {
			return strconv.FormatFloat(i, 'f', -1, 64), nil
		}
		if s, ok := val.(string); ok {
			return s, nil
		}
		return
	case TypeInt, TypeState:
		// FIX case: ES return type is unstable
		if s, ok := val.(string); ok {
			return s, nil
		}
		intVal := toInt(val)
		v = strconv.Itoa(intVal)
		return
	case TypeGroup:
		// GROUP stands form group name
		if s, ok := val.(string); ok {
			return s, nil
		}
		return
	case TypeArray, TypeJson:
		if s, ok := val.(string); ok {
			return s, nil
		}
		// fix []uint8 as json field
		if byteArr, ok := val.([]uint8); ok {
			strVal := string(byteArr)
			if _, err := json.Marshal(strVal); err == nil {
				return strVal, nil
			} else {
				logger.Error("Unable to convert json to string", "val", val, "error", err)
			}
		}
		// not []byte type
		if b, err := json.Marshal(val); err == nil {
			return string(b), nil
		} else {
			logger.Error("Unable to convert json to string", "val", val, "error", err)
		}
		return
	case TypeString, TypeText:
		if s, ok := val.(string); ok {
			return s, nil
		}
		logger.Error("Unable to convert state to string", "val", val)
		return
	case TypeTimestamp:
		// fix timestamp string as json field
		if s, ok := val.(string); ok {
			t, err := convertRFC3339(s)
			if err != nil {
				logger.Error("Unable to convert timestamp", "val", val, "s", s, "error", err)
			} else {
				return t, nil
			}
		}
		// for key and not group parameter
		if t, ok := val.(time.Time); ok {
			return Format(t), nil
		} else {
			logger.Error("Unable to convert timestamp to string", "val", val)
		}
		return
	default:
		return
	}
}

// toInt 尝试将 interface{} 类型的值转换为 int 类型，仅当值是数字类型时
func toInt(v interface{}) int {
	// 使用 type switch with assignment 来避免重复的类型断言
	switch tv := v.(type) {
	case float64:
		return int(tv)
	case float32:
		return int(tv)
	case int:
		return tv
	case int64:
		return int(tv)
	case int32:
		return int(tv)
	case int16:
		return int(tv)
	case int8:
		return int(tv)
	case uint:
		return int(tv) // #nosec G115 - type conversion is safe for typical use cases
	case uint64:
		return int(tv) // #nosec G115 - type conversion is safe for typical use cases
	case uint32:
		return int(tv)
	case uint16:
		return int(tv)
	case uint8:
		return int(tv)
	default:
		// 如果 v 不是数值类型，返回 0
		return 0
	}
}

func ParseTimeInLocation(timeStr string) (time.Time, error) {
	once.Do(initLocation)
	return time.ParseInLocation(layoutStr, timeStr, location)
}

func Format(time time.Time) string {
	return time.Format(layoutStr)
}

func convertRFC3339(rfcTimeStr string) (string, error) {
	return convertRFC3339ToCustomFormat(rfcTimeStr, layoutStr)
}

func convertRFC3339ToCustomFormat(rfcTimeStr, customFormat string) (string, error) {
	// 解析 RFC 3339 格式的时间字符串为 time.Time 对象
	parsedTime, err := time.Parse(time.RFC3339, rfcTimeStr)
	if err != nil {
		return "", err
	}
	customTimeStr := parsedTime.Format(customFormat)
	return customTimeStr, nil
}
