// Package mysql contains mysql dialect factory.
package mysql

import (
	"github.com/gantries/knife/pkg/orm"
	"github.com/gantries/knife/pkg/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	orm.RegistryDialectFactory(types.MySQL, func(properties orm.DatabaseProperties) gorm.Dialector {
		return mysql.New(mysql.Config{DriverName: properties.GetDriver(), DSN: properties.GetDSN()})
	})
}
