// Package oracle contains oracle dialect factory.
package oracle

import (
	"github.com/gantries/knife/pkg/orm"
	"github.com/gantries/knife/pkg/types"
	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
)

func init() {
	orm.RegistryDialectFactory(types.Oracle, func(properties orm.DatabaseProperties) gorm.Dialector {
		return oracle.New(oracle.Config{
			DSN:                     properties.GetDSN(),
			IgnoreCase:              false,
			NamingCaseSensitive:     true,
			VarcharSizeIsCharLength: true,
		})
	})
}
