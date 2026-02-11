// Package sqlite contains sqlite dialect factory.
package sqlite

import (
	"github.com/gantries/knife/pkg/orm"
	"github.com/gantries/knife/pkg/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	orm.RegistryDialectFactory(types.SQLite, func(properties orm.DatabaseProperties) gorm.Dialector {
		return sqlite.Open(properties.GetDSN())
	})
}
