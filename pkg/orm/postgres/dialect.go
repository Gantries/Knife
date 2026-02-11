// Package postgres contains postgres dialect factory.
package postgres

import (
	"github.com/gantries/knife/pkg/orm"
	"github.com/gantries/knife/pkg/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	orm.RegistryDialectFactory(types.Postgres, func(properties orm.DatabaseProperties) gorm.Dialector {
		return postgres.Open(properties.GetDSN())
	})
}
