package orm

import (
	"testing"
	"time"

	"github.com/gantries/knife/pkg/types"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type props struct{}

func (p props) GetDialect() types.DatabaseType    { return "unsupported" }
func (p props) GetDSN() string                    { return "" }
func (p props) GetDriver() string                 { return "" }
func (p props) GetMaxIdleConnections() int        { return 10 }
func (p props) GetMaxOpenConnections() int        { return 10 }
func (p props) GetMaxTableNameLength() int        { return 10 }
func (p props) GetConnMaxIdleTime() time.Duration { return 0 }
func (p props) GetCreateBatchSize() int           { return 10 }
func (p props) GetLogLevel() int                  { return 10 }
func (p props) ShouldPrepareStmt() bool           { return true }
func (p props) GetTablePrefix() string            { return "" }
func (p props) GetSingularTable() bool            { return true }
func (p props) GetNameReplacer() schema.Replacer  { return nil }
func (p props) GetNoLowerCase() bool              { return false }
func (p props) GetIdentifierMaxLength() int       { return 64 }
func (p props) Options() map[string]string        { return map[string]string{} }

func TestRegistryDialectFactory(t *testing.T) {
	// Test that the dialect registry works correctly
	t.Run("register and retrieve dialect", func(t *testing.T) {
		RegistryDialectFactory(types.Postgres, func(p DatabaseProperties) gorm.Dialector {
			return nil
		})
		// The factory is called when New is invoked with the registered dialect
		// Note: We can't actually test this without mocking New or causing a panic
		// The registry itself just stores the factory for later use
		// This test verifies that registration doesn't cause errors
	})
}

func TestNewPanicUnsupportedDialect(t *testing.T) {
	// Test that New panics with unsupported dialect
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unsupported dialect")
		}
	}()

	New(&props{})
}
