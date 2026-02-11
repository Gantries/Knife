package types

import (
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/gantries/knife/pkg/errors"
)

type EndpointSourceType string

const (
	ServiceEndpoint  EndpointSourceType = "Service"
	InternalEndpoint EndpointSourceType = "Internal"
)

func (y EndpointSourceType) String() string {
	return string(y)
}

func (y *EndpointSourceType) UnmarshalGQL(v interface{}) error {
	if s, ok := v.(string); ok {
		*y = EndpointSourceType(s)
		return nil
	}
	return errors.UnexpectedValueError.E(logger, "target", "endpoint-source-type", "value", v)
}

func (y EndpointSourceType) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(fmt.Sprintf(`"%s"`, y.String())))
}

func (y EndpointSourceType) Value() (driver.Value, error) {
	return string(y), nil
}

// Scan 实现了 sql.Scanner 接口
func (y *EndpointSourceType) Scan(value interface{}) error {
	// 检查值的类型是否为字符串
	switch value := value.(type) {
	case []uint8:
		*y = EndpointSourceType(value)
		return nil
	case string:
		*y = EndpointSourceType(value)
		return nil
	default:
		return errors.UnexpectedValueError.E(logger, "type", "!string/![]uint8", "value", value)
	}
}

func (c EndpointSourceType) Equal(other *EndpointSourceType) bool {
	return other != nil && c == *other
}
