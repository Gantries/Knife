package orm

import (
	"github.com/gantries/knife/pkg/types"
	"gorm.io/gorm/schema"
)

type NamingRule struct {
	strategy schema.NamingStrategy
	database types.DatabaseType
}

func (ns NamingRule) TableName(str string) (n string) {
	return ns.database.Quote(ns.strategy.TableName(str))
}

func (ns NamingRule) SchemaName(table string) string {
	return ns.database.Quote(ns.strategy.SchemaName(table))
}

func (ns NamingRule) ColumnName(table, column string) (n string) {
	return ns.database.Quote(ns.strategy.ColumnName(table, column))
}

func (ns NamingRule) JoinTableName(str string) string {
	return ns.strategy.JoinTableName(str)
}

func (ns NamingRule) RelationshipFKName(rel schema.Relationship) string {
	return ns.strategy.RelationshipFKName(rel)
}

func (ns NamingRule) CheckerName(table, column string) string {
	return ns.strategy.CheckerName(table, column)
}

func (ns NamingRule) IndexName(table, column string) string {
	return ns.strategy.IndexName(table, column)
}

func (ns NamingRule) UniqueName(table, column string) string {
	return ns.strategy.UniqueName(table, column)
}
