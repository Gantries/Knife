package orm

import (
	"testing"

	"github.com/gantries/knife/pkg/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/schema"
)

func TestNamer(t *testing.T) {
	r := NamingRule{
		schema.NamingStrategy{
			TablePrefix:         "",
			SingularTable:       true,
			IdentifierMaxLength: 64,
			NoLowerCase:         true,
			NameReplacer:        newReplacer(props{}),
		},
		types.MySQL,
	}
	assert.Equal(t, "`table`", r.TableName("`table`"), "not equal")
	assert.Equal(t, "`Table`", r.SchemaName("`table`"), "not equal")
	assert.Equal(t, "`column`", r.ColumnName("`table`", "`column`"), "not equal")
	assert.Equal(t, "`table`", r.TableName("table"), "not equal")
	assert.Equal(t, "`Table`", r.SchemaName("table"), "not equal")
	assert.Equal(t, "`column`", r.ColumnName("table", "`column`"), "not equal")
}
