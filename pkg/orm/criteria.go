package orm

import (
	"context"
	"strings"

	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/maps"
	"github.com/gantries/knife/pkg/times"
	"github.com/gantries/knife/pkg/types"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const utcTsFmt string = "2006-01-02 15:04:05"

type Criteria struct {
	query      *maps.Map[string, []any]
	orderBy    lists.List[string]
	sort       lists.List[string]
	properties DatabaseProperties
	naming     schema.Namer
	db         *Database
	table      *string
}

func Query(db *Database, table ...string) *Criteria {
	return &Criteria{
		query:      &maps.Map[string, []any]{},
		sort:       lists.List[string]{},
		orderBy:    lists.List[string]{},
		properties: db.properties,
		naming:     db.naming,
		db:         db,
		table:      lists.FirstOrDefault(table, nil),
	}
}

func (c *Criteria) put(cond string, argv []any) *Criteria {
	if !c.query.Has(cond) {
		c.sort.Add(cond)
		c.query.Put(cond, argv)
	} else {
		logger.Error("duplicate condition", "cond", cond)
	}
	return c
}

func (c *Criteria) Is(col types.ColumnName, arg any) *Criteria {
	return c.Eq(col, arg)
}

func (c *Criteria) Eq(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" = ?", []any{arg})
	return c
}

func (c *Criteria) In(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" in ?", []any{arg})
	return c
}

func (c *Criteria) Nin(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" not in ?", []any{arg})
	return c
}

func (c *Criteria) Ne(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" != ?", []any{arg})
	return c
}

func (c *Criteria) Ge(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" >= ?", []any{arg})
	return c
}

func (c *Criteria) Gt(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" > ?", []any{arg})
	return c
}

func (c *Criteria) Le(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" <= ?", []any{arg})
	return c
}

func (c *Criteria) Lt(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" < ?", []any{arg})
	return c
}

// GeUTC adds a query criterion to filter rows where the specified column's UTC timestamp is greater or equal the given timestamp.
// The given timestamp is expected to be in milliseconds since the Unix epoch. This method is specifically designed for Oracle's TIMESTAMP type.
// It formats the timestamp into a string suitable for Oracle's TIMESTAMP format and appends the condition to the query criteria.
// The function returns the current Criteria instance for chaining.
func (c *Criteria) GeUTC(col types.ColumnName, ts int64) *Criteria {
	switch c.properties.GetDialect() {
	case types.Oracle:
		c.put(c.naming.ColumnName("", col.String())+" >= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", []any{times.FormatTimestamp(ts, utcTsFmt)})
	case types.SQLServer:
		c.put(c.naming.ColumnName("", col.String())+" >= ?", []any{times.FormatTimestamp(ts, utcTsFmt)})
	}
	return c
}

// GtUTC adds a query criterion to filter rows where the specified column's UTC timestamp is greater than the given timestamp.
// The given timestamp is expected to be in milliseconds since the Unix epoch. This method is specifically designed for Oracle's TIMESTAMP type.
// It formats the timestamp into a string suitable for Oracle's TIMESTAMP format and appends the condition to the query criteria.
// The function returns the current Criteria instance for chaining.
func (c *Criteria) GtUTC(col types.ColumnName, ts int64) *Criteria {
	switch c.properties.GetDialect() {
	case types.Oracle:
		c.put(c.naming.ColumnName("", col.String())+" > TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", []any{times.FormatTimestamp(ts, utcTsFmt)})
	case types.SQLServer:
		c.put(c.naming.ColumnName("", col.String())+" > ?", []any{times.FormatTimestamp(ts, utcTsFmt)})
	}
	return c
}

// LeUTC adds a query criterion to filter rows where the specified column's UTC timestamp is less or equal the given timestamp.
// The given timestamp is expected to be in milliseconds since the Unix epoch. This method is specifically designed for Oracle's TIMESTAMP type.
// It formats the timestamp into a string suitable for Oracle's TIMESTAMP format and appends the condition to the query criteria.
// The function returns the current Criteria instance for chaining.
func (c *Criteria) LeUTC(col types.ColumnName, ts int64) *Criteria {
	switch c.properties.GetDialect() {
	case types.Oracle:
		c.put(c.naming.ColumnName("", col.String())+" <= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", []any{times.FormatTimestamp(ts, utcTsFmt)})
	case types.SQLServer:
		c.put(c.naming.ColumnName("", col.String())+" <= ?", []any{times.FormatTimestamp(ts, utcTsFmt)})
	}
	return c
}

// LtUTC adds a query criterion to filter rows where the specified column's UTC timestamp is less than the given timestamp.
// The given timestamp is expected to be in milliseconds since the Unix epoch. This method is specifically designed for Oracle's TIMESTAMP type.
// It formats the timestamp into a string suitable for Oracle's TIMESTAMP format and appends the condition to the query criteria.
// The function returns the current Criteria instance for chaining.
func (c *Criteria) LtUTC(col types.ColumnName, ts int64) *Criteria {
	switch c.properties.GetDialect() {
	case types.Oracle:
		c.put(c.naming.ColumnName("", col.String())+" < TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", []any{times.FormatTimestamp(ts, utcTsFmt)})
	case types.SQLServer:
		c.put(c.naming.ColumnName("", col.String())+" < ?", []any{times.FormatTimestamp(ts, utcTsFmt)})
	}
	return c
}

func (c *Criteria) Between(col types.ColumnName, min, max any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" between (?, ?)", []any{min, max})
	return c
}

func (c *Criteria) Like(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" like ?", []any{arg})
	return c
}

func (c *Criteria) Unlike(col types.ColumnName, arg any) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" not like ?", []any{arg})
	return c
}

func (c *Criteria) Where() maps.Map[string, []any] {
	return *c.query
}

func (c *Criteria) Build() *gorm.DB {
	return c.BuildWithContext(context.Background())
}

func (c *Criteria) BuildWithTx(tx *gorm.DB) *gorm.DB {
	db := tx
	for _, k := range c.sort {
		if (*c.query)[k] == nil {
			db = db.Where(k)
		} else {
			db = db.Where(k, (*c.query)[k]...)
		}
	}
	for _, item := range c.orderBy {
		db = db.Order(item)
	}
	return db
}

// BuildWithTxIfPossible is a method that conditionally constructs a database query,
// either within the context of a transaction or as a regular DML operation.
// It checks if the provided transaction context (tx) is nil. If tx is nil, it
// indicates that there is no ongoing transaction, and the method will proceed with
// a non-transactional DML operation by calling the Build method. Otherwise, if
// tx is not nil, it implies that the operation should be part of an existing
// transaction, and the method will use the BuildWithTx method to continue within
// the transactional context.
//
// Parameters:
//
//	tx *gorm.DB: The transactional database connection. If not nil, the query will be executed within a transaction.
//
// Returns:
//
//	*gorm.DB: The constructed database query, either transactional or non-transactional.
func (c *Criteria) BuildWithTxIfPossible(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return c.Build()
	} else {
		return c.BuildWithTx(tx)
	}
}

// BuildWithTxAndChangeTable is a method that extends the functionality of the Criteria
// struct to support transactional operations across multiple tables. It takes a
// transactional database connection and changes the table context for the operation.
// This method is particularly useful when you need to perform actions that involve
// multiple tables within the scope of a single transaction.
// The function first builds the database connection with the given transaction using
// the BuildWithTx method. It then checks if the 'table' field in the Criteria struct
// is not nil. If 'table' is not nil, indicating that a specific table has been set,
// the method changes the table context to the one specified in the Criteria struct's
// 'table' field. If 'table' is nil, no table switching is performed, and the database
// connection remains in its original context.
//
// Parameters:
//
//	tx *gorm.DB: The transactional database connection to use for the operation.
//
// Returns:
//
//	*gorm.DB: The database connection, potentially with the table context changed if 'table' is not nil.
func (c *Criteria) BuildWithTxAndChangeTable(tx *gorm.DB) *gorm.DB {
	db := c.BuildWithTx(tx)
	if c.table != nil {
		db = db.Table(*c.table)
	}
	return db
}

// BuildWithTxIfPossibleAndChangeTable checks if a transaction is provided and either continues within the transactional context
// using BuildWithTx or performs a non-transactional operation with Build. It also sets the table context if one is specified.
// See BuildWithTxIfPossible and BuildWithTxAndChangeTable for more details on transactional and table handling.
func (c *Criteria) BuildWithTxIfPossibleAndChangeTable(tx *gorm.DB) *gorm.DB {
	var db *gorm.DB
	if tx == nil {
		db = c.Build()
	} else {
		db = c.BuildWithTx(tx)
	}
	if c.table != nil {
		db = db.Table(*c.table)
	}
	return db
}

func (c *Criteria) BuildWithContext(ctx context.Context) *gorm.DB {
	if c.table != nil {
		return c.BuildWithTx(c.db.DB().WithContext(ctx).Table(*c.table))
	}
	return c.BuildWithTx(c.db.DB().WithContext(ctx))
}

func (c *Criteria) BuildQuery(sql *string) (statement string, args []any) {
	builder := strings.Builder{}
	builder.WriteString(*sql)
	args = make([]any, 0, c.query.Length()*2)
	if c.query.Length() > 0 {
		builder.WriteString(" where ")
		notFirst := false
		for _, w := range c.sort {
			if notFirst {
				builder.WriteString(" and ")
			} else {
				notFirst = true
			}
			builder.WriteString(w)
			if (*c.query)[w] != nil {
				args = append(args, (*c.query)[w]...)
			}
		}
	}

	c.orderBy.ForRest(
		1,
		func(order string) {
			builder.WriteString(" order by ")
			builder.WriteString(order)
		},
		func(order string) {
			builder.WriteString(", ")
			builder.WriteString(order)
		},
	)

	statement = builder.String()
	return
}

func (c *Criteria) Exec(ctxt context.Context, sql *string, print bool) *gorm.DB {
	statement, args := c.BuildQuery(sql)
	if print {
		logger.Info("Built sql", "sql", statement)
	}
	if c.table != nil {
		return c.db.DB().Table(*c.table).WithContext(ctxt).Raw(statement, args...)
	} else {
		return c.db.DB().WithContext(ctxt).Raw(statement, args...)
	}
}

func (c *Criteria) Asc(cols ...types.ColumnName) *Criteria {
	s := strings.Join(*lists.For(&cols, func(col types.ColumnName) string {
		return c.naming.ColumnName("", col.String())
	}), ", ")
	c.orderBy.Add(s + " asc")
	return c
}

func (c *Criteria) Desc(cols ...types.ColumnName) *Criteria {
	s := strings.Join(*lists.For(&cols, func(col types.ColumnName) string {
		return c.naming.ColumnName("", col.String())
	}), ", ")
	c.orderBy.Add(s + " desc")
	return c
}

func (c *Criteria) Null(col types.ColumnName) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" is null", nil)
	return c
}

func (c *Criteria) NotNull(col types.ColumnName) *Criteria {
	c.put(c.naming.ColumnName("", col.String())+" is not null", nil)
	return c
}
