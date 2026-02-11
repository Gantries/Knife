package types

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strings"

	"github.com/gantries/knife/pkg/errors"
)

type ColumnName string

const GroupSeparator = "/"
const ColumnSeparator = "_"

const (
	ColumnApproveReason       ColumnName = "approve_reason"
	ColumnBpmNo               ColumnName = "bpm_no"
	ColumnBpmStatus           ColumnName = "bpm_status"
	ColumnCreatorId           ColumnName = "creator_id"
	ColumnCreatorDisplayName  ColumnName = "creator_display_name"
	ColumnCreateTime          ColumnName = "create_time"
	ColumnDeleteTime          ColumnName = "delete_time"
	ColumnDeleteFlag          ColumnName = "delete_flag"
	ColumnIdentifier          ColumnName = "id"
	ColumnMatchTypes          ColumnName = "match_types"
	ColumnModifierId          ColumnName = "modifier_id"
	ColumnModifierDisplayName ColumnName = "modifier_display_name"
	ColumnModifiedBy          ColumnName = "modified_by"
	ColumnModify              ColumnName = "modify"
	ColumnModifyTime          ColumnName = "modify_time"
	ColumnPublishedFlag       ColumnName = "published_flag"
	ColumnStatus              ColumnName = "status"
	ColumnUpdateTime          ColumnName = "update_time"
	ColumnValidFlag           ColumnName = "valid_flag"
	ColumnVersion             ColumnName = "version"
	ColumnWeight              ColumnName = "weight"
	ColumnLastReleasedTime    ColumnName = "latest_released_time"
)

func (c ColumnName) String() string {
	return string(c)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (c *ColumnName) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if ok {
		*c = ColumnName(str)
	} else {
		*c = v.(ColumnName)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (c ColumnName) MarshalGQL(w io.Writer) {
	_, err := w.Write([]byte(fmt.Sprintf(`"%s"`, c.String())))
	if err != nil {
		return
	}
}

func (c ColumnName) Dup() *string {
	p := string(c)
	return &p
}

func (c ColumnName) Group(argv ...string) string {
	return strings.Join(append(argv, c.String()), GroupSeparator)
}

func (c ColumnName) Append(argv ...string) ColumnName {
	return ColumnName(strings.Join(append([]string{c.String()}, argv...), ""))
}

func GroupToColumnName(group string) ColumnName {
	return ColumnName(strings.ReplaceAll(group, GroupSeparator, ColumnSeparator))
}

func (c ColumnName) Value() (driver.Value, error) {
	return string(c), nil
}

// Scan 实现了 sql.Scanner 接口
func (c *ColumnName) Scan(value interface{}) error {
	switch value := value.(type) {
	case []uint8:
		*c = ColumnName(value)
		return nil
	case string:
		*c = ColumnName(value)
		return nil
	}
	return errors.UnexpectedValueError.E(logger, "type", "!string/![]uint8", "value", value)
}

func (c ColumnName) Equal(other *ColumnName) bool {
	return other != nil && c == *other
}
