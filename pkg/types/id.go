// Package types contains all types definition for gqlgin/gorm
package types

import (
	"fmt"
	"io"
	"strconv"

	"github.com/gantries/knife/pkg/lang"
)

// ID type alias
type ID int64

// Id used for backward compatibility
type Id = ID

// Timestamp used for timestamp(int64)
type Timestamp = ID

// Value returns int64 representive
func (y *ID) Value() int64 {
	return int64(*y)
}

// Ptr returns pointer to the actual value
func (y *ID) Ptr() *int64 {
	v := y.Value()
	return &v
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (y *ID) UnmarshalGQL(v interface{}) error {
	ptr := lang.ToLong(v, func(i *int64) *ID { t := ID(*i); return &t }, nil)
	if ptr != nil {
		*y = *ptr
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (y ID) MarshalGQL(w io.Writer) {
	_, err := w.Write([]byte(fmt.Sprintf(`"%d"`, int64(y))))
	if err != nil {
		return
	}
}

// String converts ID instance to string expression
func String(id ID) string {
	return strconv.FormatInt(int64(id), 10)
}
