package orm

import (
	"testing"

	"github.com/gantries/knife/pkg/lang"
	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/schema"
)

func TestQuery(t *testing.T) {
	//var props = properties{}
	var db = Database{
		naming: NamingRule{
			strategy: schema.NamingStrategy{
				TablePrefix:         "",
				SingularTable:       true,
				IdentifierMaxLength: 256,
				NoLowerCase:         false,
			},
			database: types.MySQL,
		},
	}
	q := db.Query("table")
	assert.True(t, (q.Is(types.ColumnIdentifier, 123).Where())["`id` = ?"][0] == 123)
	assert.True(t, ((q.In(types.ColumnIdentifier, *lists.Of(123, 456)).Where())["`id` in ?"][0].([]int))[0] == 123)
	assert.True(t, (q.Ne(types.ColumnIdentifier, 123).Where())["`id` != ?"][0] == 123)
	assert.True(t, (q.Ge(types.ColumnIdentifier, 123).Where())["`id` >= ?"][0] == 123)
	assert.True(t, (q.Gt(types.ColumnIdentifier, 123).Where())["`id` > ?"][0] == 123)
	assert.True(t, (q.Le(types.ColumnIdentifier, 123).Where())["`id` <= ?"][0] == 123)
	assert.True(t, (q.Null(types.ColumnIdentifier).Where())["`id` is null"] == nil)
	assert.True(t, (q.NotNull(types.ColumnIdentifier).Where())["`id` is not null"] == nil)
	assert.True(t, (q.Lt(types.ColumnIdentifier, 123).Where())["`id` < ?"][0] == 123)
	assert.True(t, (q.Between(types.ColumnIdentifier, 123, 456).Where())["`id` between (?, ?)"][0] == 123)
	// Note: GtUTC and LtUTC require properties to be set, so we skip those for this unit test
	q.Asc("asc1", "asc2")
	q.Desc("desc1", "desc2", "desc3")
	s, _ := q.BuildQuery(lang.Dup(""))
	assert.Equal(t, " where `id` = ? and `id` in ? and `id` != ? and `id` >= ? and `id` > ? and `id` <= ? and `id` is null and `id` is not null and `id` < ? and `id` between (?, ?) order by `asc1`, `asc2` asc, `desc1`, `desc2`, `desc3` desc", s)
}
