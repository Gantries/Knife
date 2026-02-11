package orm

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/lang"
	"github.com/gantries/knife/pkg/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type transactionKeyType struct{}

var transactionKey = transactionKeyType{}

type transaction struct {
	tx       *gorm.DB
	database string
}

type Database struct {
	db         *gorm.DB
	raw        *sql.DB
	properties DatabaseProperties
	naming     schema.Namer
	database   string
}

func (d *Database) DB() *gorm.DB {
	return d.db
}

func (d *Database) EscapeCharacters() (string, string) {
	return DatabaseEscapeCharacters(d.properties.GetDialect())
}

func (d *Database) Query(table ...string) *Criteria {
	return Query(d, table...)
}

func (d *Database) Table(table string) *gorm.DB {
	return d.db.Table(table)
}

func (d *Database) TableWithContext(ctxt context.Context, table string) *gorm.DB {
	return d.db.Table(table).WithContext(ctxt)
}

func (d *Database) Escape(c types.ColumnName) string {
	return d.naming.ColumnName("", c.String())
}

func DatabaseEscapeCharacters(typ types.DatabaseType) (string, string) {
	if v, ok := types.DatabaseEscapeCharacters[typ]; ok {
		return v.Left, v.Right
	}
	return "", ""
}

func (d *Database) Convert(s string, t types.ParameterType, tr *i18n.Localizer, log *slog.Logger) (v interface{}, err error) {
	if len(s) <= 0 {
		switch t {
		default:
			return nil, nil
		}
	}
	switch d.properties.GetDialect() {
	case types.MySQL, types.Postgres, types.Oracle, types.SQLite, types.DB2:
		switch t {
		case types.TypeBool:
			return lang.Ternary(trueEquivalents.Has(strings.ToLower(s)), 1, 0), nil
		case types.TypeDouble:
			return strconv.ParseFloat(s, 64)
		case types.TypeIdentity, types.TypeTimestamp:
			return strconv.ParseInt(s, 10, 64)
		case types.TypeInt:
			return strconv.ParseInt(s, 10, 32)
		case types.TypeGroup, types.TypeJson, types.TypeString, types.TypeText, types.TypeArray:
			return s, nil
		case types.TypeState:
			return strconv.ParseInt(s, 10, 8)
		}
	//	sqlserver use varbinary save json and array
	case types.SQLServer:
		switch t {
		case types.TypeBool:
			return lang.Ternary(trueEquivalents.Has(strings.ToLower(s)), 1, 0), nil
		case types.TypeDouble:
			return strconv.ParseFloat(s, 64)
		case types.TypeIdentity:
			return strconv.ParseInt(s, 10, 64)
		case types.TypeInt:
			return strconv.ParseInt(s, 10, 32)
		case types.TypeGroup, types.TypeString, types.TypeText:
			return s, nil
		case types.TypeJson, types.TypeArray:
			return []byte(s), nil
		case types.TypeState:
			return strconv.ParseInt(s, 10, 8)
		case types.TypeTimestamp:
			time, err := types.ParseTimeInLocation(s)
			if err != nil {
				logger.Error("cannot parse time string to time.Time", "s", s, "error", err)
				return nil, err
			}
			return time, nil
		}
	}
	return nil, errors.UnexpectedTypeError.LocalE(tr, log, "type", t)
}

// Tx can be used to execute statement in an existing transaction, or in a
// newly created transaction.
func (d *Database) Tx(ctx context.Context, fc func(c context.Context, tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	if tx := d.OptionalTx(ctx); tx != nil {
		return fc(ctx, tx)
	}
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fc(context.WithValue(ctx, transactionKey, transaction{tx, d.database}), tx)
	}, opts...)
}

// OptionalTx can be used to retrieve a [gorm.DB] pointer bound to given [context.Context].
func (d *Database) OptionalTx(ctx context.Context) *gorm.DB {
	if v := ctx.Value(transactionKey); v != nil {
		if tr, ok := v.(transaction); ok {
			if tr.database != d.database {
				logger.Warn("Using a transaction from a different database may cause problem", "db", tr.database, "current-db", d.database)
			}
			return tr.tx
		}
	}
	return nil
}
