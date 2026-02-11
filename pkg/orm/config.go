package orm

import (
	"time"

	"github.com/gantries/knife/pkg/types"
	"gorm.io/gorm/schema"
)

type DatabaseProperties interface {
	GetDSN() string
	GetDialect() types.DatabaseType
	GetDriver() string
	GetMaxIdleConnections() int
	GetMaxOpenConnections() int
	GetMaxTableNameLength() int
	GetConnMaxIdleTime() time.Duration
	GetCreateBatchSize() int
	GetLogLevel() int // see gorm.io/gorm/logger.LogLevel
	ShouldPrepareStmt() bool
	GetTablePrefix() string
	GetSingularTable() bool
	GetNameReplacer() schema.Replacer
	GetNoLowerCase() bool
	GetIdentifierMaxLength() int
	Options() map[string]string
}
