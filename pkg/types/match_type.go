package types

import (
	"fmt"
	"io"

	"github.com/gantries/knife/pkg/errors"
)

type MatchType string

const (
	MatchAny    MatchType = "any"
	MatchEqual  MatchType = "equal"
	MatchGroup  MatchType = "group"
	MatchPrefix MatchType = "prefix"
)

func (y MatchType) String() string {
	return string(y)
}

func (y *MatchType) UnmarshalGQL(v interface{}) error {
	if s, ok := v.(string); ok {
		switch MatchType(s) {
		case MatchEqual:
			*y = MatchType(s)
			return nil
		case MatchPrefix:
			*y = MatchType(s)
			return nil
		case MatchAny:
			*y = MatchType(s)
			return nil
		case MatchGroup:
			*y = MatchType(s)
			return nil
		default:
		}
	}
	return errors.UnexpectedValueError.E(logger, "target", "match-type", "value", v)
}

func (y MatchType) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(fmt.Sprintf(`"%s"`, y.String())))
}

// Scan 实现了 sql.Scanner 接口
func (c *MatchType) Scan(value interface{}) error {
	switch value := value.(type) {
	case []uint8:
		*c = MatchType(value)
		return nil
	case string:
		*c = MatchType(value)
		return nil
	}
	return errors.UnexpectedValueError.E(logger, "type", "!string/![]uint8", "value", value)
}

func (c MatchType) Equal(other *MatchType) bool {
	return other != nil && c == *other
}
