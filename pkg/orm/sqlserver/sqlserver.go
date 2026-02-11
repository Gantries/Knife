// Package sqlserver contains sqlserver dialect factory.
package sqlserver

import (
	"github.com/gantries/knife/pkg/orm"
	"github.com/gantries/knife/pkg/types"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func init() {
	orm.RegistryDialectFactory(types.SQLServer, func(properties orm.DatabaseProperties) gorm.Dialector {
		return sqlserver.Open(properties.GetDSN())
	})
}
