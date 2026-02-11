package orm

import (
	"fmt"
	"strings"

	"github.com/gantries/knife/pkg/lang"
	"github.com/gantries/knife/pkg/log"
	"github.com/gantries/knife/pkg/types"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	At       string = "@"
	Question string = "?"
)

const IdentifierMaxLength = 64

var logger = log.New("knife/orm")

var dialectorFactories = map[types.DatabaseType]func(DatabaseProperties) gorm.Dialector{}

func RegistryDialectFactory(t types.DatabaseType, factory func(DatabaseProperties) gorm.Dialector) {
	dialectorFactories[t] = factory
}

func New(properties DatabaseProperties) *Database {
	var dialect gorm.Dialector

	var gen = NamingRule{
		schema.NamingStrategy{
			TablePrefix:         properties.GetTablePrefix(),
			SingularTable:       properties.GetSingularTable(),
			IdentifierMaxLength: properties.GetIdentifierMaxLength(),
			NoLowerCase:         properties.GetNoLowerCase(),
			NameReplacer:        properties.GetNameReplacer(),
		},
		properties.GetDialect(),
	}

	if gen.strategy.NameReplacer == nil {
		gen.strategy.NameReplacer = newReplacer(properties)
	}

	if f, ok := dialectorFactories[properties.GetDialect()]; ok {
		dialect = f(properties)
	} else {
		panic(fmt.Errorf("unsupported database dialect: %s", properties.GetDialect()))
	}

	db, err := gorm.Open(dialect, &gorm.Config{
		SkipDefaultTransaction:   false,
		PrepareStmt:              properties.ShouldPrepareStmt(),
		NamingStrategy:           gen,
		DryRun:                   false,
		DisableAutomaticPing:     false,
		DisableNestedTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	db.Logger = db.Logger.LogMode(gormlog.LogLevel(properties.GetLogLevel()))
	raw, err := db.DB()
	if err != nil {
		panic(err)
	}
	raw.SetMaxIdleConns(properties.GetMaxIdleConnections())
	raw.SetMaxOpenConns(properties.GetMaxOpenConnections())
	raw.SetConnMaxIdleTime(properties.GetConnMaxIdleTime())

	dbid := properties.GetDSN()
	beg, end := strings.Index(dbid, At), strings.LastIndex(dbid, Question)
	beg, end = lang.Ternary(beg >= 0, beg+1, 0), lang.Ternary(end > 0, end, len(dbid))
	return &Database{db, raw, properties, gen, dbid[beg:end]}
}
