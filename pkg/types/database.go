package types

import (
	"database/sql/driver"
	"io"

	"github.com/gantries/knife/pkg/errors"
)

type DatabaseType string

const (
	MySQL     DatabaseType = "mysql"     // gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local
	Postgres  DatabaseType = "postgres"  // host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	Oracle    DatabaseType = "oracle"    // oracle://user:password@127.0.0.1:1521/service
	SQLite    DatabaseType = "sqlite"    // test.db
	SQLServer DatabaseType = "sqlserver" // sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm
	DB2       DatabaseType = "db2"
)

var DatabaseEscapeCharacters = map[DatabaseType]struct{ Left, Right string }{
	MySQL:     {"`", "`"},
	Postgres:  {"\"", "\""},
	Oracle:    {"\"", "\""},
	SQLite:    {"\"", "\""},
	DB2:       {"\"", "\""},
	SQLServer: {"[", "]"},
}

func (c DatabaseType) String() string {
	return string(c)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (c *DatabaseType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		*c = DatabaseType(str)
	} else {
		*c = v.(DatabaseType)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (c DatabaseType) MarshalGQL(w io.Writer) {
	_, err := w.Write([]byte(c.String()))
	if err != nil {
		return
	}
}

func (c DatabaseType) Value() (driver.Value, error) {
	return string(c), nil
}

// Scan 实现了 sql.Scanner 接口
func (c *DatabaseType) Scan(value interface{}) error {
	switch value := value.(type) {
	case []uint8:
		*c = DatabaseType(value)
		return nil
	case string:
		*c = DatabaseType(value)
		return nil
	}
	return errors.UnexpectedValueError.E(logger, "type", "!string/![]uint8", "value", value)
}

func (c DatabaseType) Equal(other *DatabaseType) bool {
	return other != nil && c == *other
}

func (c DatabaseType) Quote(i string) (o string) {
	v, ok := DatabaseEscapeCharacters[c]
	if !ok {
		return i
	}

	o = i
	if l := len(o); l > 0 {
		nstart, nstop := o[0:1] != v.Left, o[l-1:l] != v.Right
		if nstart {
			o = v.Left + o
		}
		if nstop {
			o = o + v.Right
		}
	}
	return
}
