# knife API Documentation for AI Agents

This document provides comprehensive API reference for AI agents to understand and use the knife Go utility library effectively.

## Table of Contents

- [Overview](#overview)
- [Package Categories](#package-categories)
- [Common Patterns](#common-patterns)
- [Best Practices](#best-practices)

---

## Overview

**knife** is a comprehensive Go utility library providing:

- Database ORM with multi-database support
- Caching layer with Redis
- Serialization (JSON/YAML)
- Expression evaluation
- OpenTelemetry integration
- Authentication & authorization
- Elasticsearch client
- Internationalization (i18n)
- Common utilities (maps, lists, time, etc.)

**Module Path:** `github.com/gantries/knife`

**Import Example:**

```go
import "github.com/gantries/knife/pkg/orm"
```

---

## Package Categories

### 1. Database ORM (`pkg/orm/`)

Unified database interface supporting MySQL, PostgreSQL, Oracle, SQLite, SQL Server.

**Key Types:**

```go
type Database struct {
    db         *gorm.DB
    raw        *sql.DB
    properties DatabaseProperties
    naming     schema.Namer
    database   string
}

type Criteria struct {
    query      *maps.Map[string, []any]
    orderBy    lists.List[string]
    sort       lists.List[string]
    distinct   bool
    forUpdate  bool
}
```

**Database Properties Interface:**

```go
type DatabaseProperties interface {
    GetDSN() string
    GetDialect() types.DatabaseType
    GetTablePrefix() string
    GetSingularTable() bool
    GetNoLowerCase() bool
    GetIdentifierMaxLength() int
    GetNameReplacer() strings.Replacer
    GetLogLevel() int
    GetMaxIdleConnections() int
    GetMaxOpenConnections() int
    GetConnMaxIdleTime() time.Duration
    GetConnMaxLifeTime() time.Duration
    ShouldPrepareStmt() bool
}
```

**Key Functions:**

```go
// Initialize database
func New(properties DatabaseProperties) *Database

// Query builder
func (d *Database) Query(name string) *Criteria

// Transaction wrapper
func (d *Database) Tx(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error

// Access underlying GORM
func (d *Database) DB() *gorm.DB
func (d *Database) Raw() *sql.DB
```

**Criteria Builder Methods:**

```go
// Where conditions
func (c *Criteria) Eq(column types.ColumnName, value any) *Criteria
func (c *Criteria) Ne(column types.ColumnName, value any) *Criteria
func (c *Criteria) Gt(column types.ColumnName, value any) *Criteria
func (c *Criteria) Gte(column types.ColumnName, value any) *Criteria
func (c *Criteria) Lt(column types.ColumnName, value any) *Criteria
func (c *Criteria) Lte(column types.ColumnName, value any) *Criteria
func (c *Criteria) Like(column types.ColumnName, value any) *Criteria
func (c *Criteria) In(column types.ColumnName, values ...any) *Criteria
func (c *Criteria) NotIn(column types.ColumnName, values ...any) *Criteria
func (c *Criteria) IsNull(column types.ColumnName) *Criteria
func (c *Criteria) IsNotNull(column types.ColumnName) *Criteria
func (c *Criteria) Between(column types.ColumnName, min, max any) *Criteria

// Sorting
func (c *Criteria) Asc(column types.ColumnName) *Criteria
func (c *Criteria) Desc(column types.ColumnName) *Criteria

// Modifiers
func (c *Criteria) Distinct() *Criteria
func (c *Criteria) ForUpdate() *Criteria
func (c *Criteria) Limit(limit int) *Criteria
func (c *Criteria) Offset(offset int) *Criteria

// Execution
func (c *Criteria) Build() *gorm.DB
func (c *Criteria) Count() (int64, error)
func (c *Criteria) Exists() (bool, error)
func (c *Criteria) First() *gorm.DB
func (c *Criteria) Find() *gorm.DB
func (c *Criteria) Scan(dest any) *gorm.DB
```

**Supported Databases:**

```go
const (
    MySQL      DatabaseType = "mysql"
    Postgres   DatabaseType = "postgres"
    Oracle     DatabaseType = "oracle"
    SQLite     DatabaseType = "sqlite"
    SQLServer  DatabaseType = "sqlserver"
    DB2        DatabaseType = "db2"
)
```

**Usage Example:**

```go
// Initialize
db := orm.New(&MyProperties{
    Dialect:      types.Postgres,
    DSN:          "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable",
    TablePrefix:  "app_",
    SingularTable: true,
})

// Query with type-safe criteria
users := db.Query("users").
    Eq(types.ColumnName("status"), "active").
    Gt(types.ColumnName("age"), 18).
    Asc(types.ColumnName("name")).
    Limit(10).
    Find(&[]User{})

// Transaction
err := db.Tx(ctx, func(ctx context.Context, tx *gorm.DB) error {
    return tx.Create(&User{Name: "John"}).Error
})
```

---

### 2. Caching Layer (`pkg/cache/`)

Generic cache abstraction with Redis implementation and HA support.

**Key Types:**

```go
// Cache interface
type Cache interface {
    Ping(ctx context.Context) error

    // String operations
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
    Del(ctx context.Context, keys ...string) (int64, error)
    Exists(ctx context.Context, keys ...string) (int64, error)

    // List operations
    Push(ctx context.Context, key string, values ...interface{}) error
    Pop(ctx context.Context, key string) (string, error)
    LPop(ctx context.Context, key string) (string, error)
    RPush(ctx context.Context, key string, fields ...string) (int64, error)

    // Hash operations
    HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
    HGet(ctx context.Context, key, field string) (string, error)
    HGetAll(ctx context.Context, key string) (map[string]string, error)
    HDel(ctx context.Context, key string, fields ...string) (int64, error)

    // Counters
    Count(ctx context.Context, key string) (int64, error)
    Incr(ctx context.Context, key string) (int64, error)
    IncrBy(ctx context.Context, key string, value int64) (int64, error)
}

// Redis cache implementation
type RedisCache struct {
    client      *redis.Client
    properties  *Properties
}

// Configuration
type Properties struct {
    Type       Type
    Addresses  lists.List[string]
    Database   int
    Credential Credential
    Pool       Pool
    Username   string
    Password   string
    Addr       string
}

type Credential struct {
    Username string
    Password string
}

type Pool struct {
    Size                  int
    MaxRetries            int
    MinIdleConnections    int
    MaxIdleConnections    int
    MaxActiveConnections  int
    MaxConnectionIdleTime time.Duration
    MaxConnectionLifeTime time.Duration
}
```

**Key Functions:**

```go
// Create cache instance
func New(ctx context.Context, props *Properties) (Cache, *national.Message)
func NewRedis(props *Properties) (Cache, error)
```

**Usage Example:**

```go
// Create Redis cache
cache, err := cache.New(ctx, &cache.Properties{
    Type:      cache.Redis,
    Addresses: lists.Of[string]("localhost:6379"),
    Database:  0,
    Pool: cache.Pool{
        MaxIdleConnections:   10,
        MaxActiveConnections: 30,
    },
})

// String operations
cache.Set(ctx, "user:1", "data", time.Hour)
value, err := cache.Get(ctx, "user:1")

// Hash operations
cache.HSet(ctx, "user:1", "name", "John", "age", "30")
name, err := cache.HGet(ctx, "user:1", "name")
all, err := cache.HGetAll(ctx, "user:1")

// List operations
cache.RPush(ctx, "queue", "job1", "job2")
job, err := cache.LPop(ctx, "queue")
```

---

### 3. Serialization (`pkg/serde/`)

Generic JSON and YAML serialization with type safety.

**Key Functions:**

```go
// JSON operations
func Serialize[T any](v T) (buf []byte, err error)
func Deserialize[T any](b []byte) (t *T, err error)
func DeserializeArray[T any](b []byte) (a []T, err error)

// YAML operations
func SerializeYAML[T any](v T) (buf []byte, err error)
func DeserializeYAML[T any](b []byte) (t *T, err error)
func DeserializeYAMLArray[T any](b []byte) (a []T, err error)

// Validation
func IsJSON(s string) bool
func IsJSONArray(s string) bool
```

**Usage Example:**

```go
// JSON with generics
user := User{Name: "John", Age: 30}
data, err := serde.Serialize(user)
user, err := serde.Deserialize[User](data)

// Array handling
users := []User{{Name: "John"}, {Name: "Jane"}}
data, err := serde.Serialize(users)
result, err := serde.DeserializeArray[User](data)

// YAML
config, err := serde.DeserializeYAML[Config](yamlBytes)
```

---

### 4. Expression Evaluation (`pkg/eval/`)

Safe runtime expression evaluation using expr library.

**Key Functions:**

```go
// Compile and evaluate
func Compile(program string) (*vm.Program, error)

func Evaluate(tr *i18n.Localizer, tpl *string, vars maps.Map[string, interface{}],
    envs ...maps.Map[string, interface{}]) (*string, error)

// Template evaluation
func EvaluateTemplate(tr *i18n.Localizer, tpl *string, vars maps.Map[string, interface{}],
    envs ...maps.Map[string, interface{}]) (*string, error)

// Built-in functions available in expressions:
// - fmt(format, ...args): String formatting
// - log: Logger instance
// - s2i(...pairs): Create string map
// - i2i(...pairs): Create interface map
// - arr(...items): Create array
// - p2i32, p2i, p2s, p2a, p2f32, p2f, p2b: Default value coercions
```

**Usage Example:**

```go
// Simple evaluation
result, err := eval.Evaluate(
    national.Tr(ctx),
    &"user.Age > 18 && user.Active == true",
    maps.Of[string, interface{}]("user", map[string]interface{}{
        "Age": 25, "Active": true,
    }),
)

// With multiple environments
vars := maps.Of[string, interface{}]("x", 10)
env1 := maps.Of[string, interface{}]("y", 20)
result, err := eval.Evaluate(nil, &"x + y", vars, env1)

// Template with i18n
template := "Hello {{.name}}, you have {{.count}} messages"
result, err := eval.EvaluateTemplate(
    national.Tr(ctx),
    &template,
    maps.Of[string, interface{}]("name", "John", "count", 5),
)
```

---

### 5. OpenTelemetry (`pkg/tel/`)

OpenTelemetry integration for tracing, metrics, and logging.

**Key Types:**

```go
// Configuration
type Remote struct {
    Protocol EndpointProtocol
    Endpoint string
}

type EndpointProtocol string

const (
    HTTP  EndpointProtocol = "http"
    GRPC  EndpointProtocol = "grpc"
)

type OpenTelemetry struct {
    Trace   Remote
    Metric  Remote
    Log     Remote
    Disable bool
}

// Meter and tracer interfaces
type SimpleMeter interface {
    // Counters
    NewCounter(name string, opts ...metric.InstrumentOption) SimpleCounter

    // Histograms
    NewHistogram(name string, opts ...metric.InstrumentOption) SimpleHistogram
}

type SimpleCounter interface {
    Add(ctx context.Context, incr int64, options ...metric.AddOption)
}

type SimpleHistogram interface {
    Record(ctx context.Context, incr float64, options ...metric.RecordOption)
}
```

**Key Functions:**

```go
// Setup SDK
func SetupOTelSDK(ctx context.Context, name string, config OpenTelemetry)
    (shutdown func(context.Context) error, err error)

// Individual components
func SetupTracer(name string, remote Remote) (*trace.TracerProvider, error)
func SetupMeter(name string, remote Remote) (*metric.MeterProvider, error)
func SetupLoggersCreated(p log.LoggerProvider)

// Shutdown
func ShutdownTracer(ctx context.Context, tp *trace.TracerProvider) error
func ShutdownMeterProvider(ctx context.Context, mp *metric.MeterProvider) error

// Logging
func Logger(name string) *slog.Logger

// Built-in attributes
func BuiltinAttributes() lists.List[attribute.KeyValue]
func BuiltinAttributeStrings() lists.List[string]
func BuiltinAttributeFlatStrings() lists.List[string]
func Fingerprint() (hostname string, ip, mac []string, ns string)
```

**Usage Example:**

```go
// Setup complete OTel pipeline
config := tel.OpenTelemetry{
    Trace: tel.Remote{Protocol: tel.GRPC, Endpoint: "localhost:4317"},
    Metric: tel.Remote{Protocol: tel.HTTP, Endpoint: "localhost:4318"},
    Log: tel.Remote{Protocol: tel.GRPC, Endpoint: "localhost:4319"},
}

shutdown, err := tel.SetupOTelSDK(ctx, "my-service", config)
defer shutdown(ctx)

// Use logger
logger := tel.Logger("my-app")
logger.Info("Processing request", "user_id", 123)

// Use built-in attributes
attrs := tel.BuiltinAttributes()
```

---

### 6. Authentication (`pkg/auth/`)

Authentication and authorization utilities.

**Key Types:**

```go
type Identity struct {
    Email             string
    Name              string
    UserName          string
    Raw               string
    authenticated     bool
    authenticatedKeys map[string]bool
}

type contextKeyType string

const (
    HeaderIdentity = contextKeyType("x-userinfo")
    KeyIdentity    = contextKeyType("authorization")
)
```

**Key Functions:**

```go
// Create identity
func NewIdentity(email, name, username, raw string) *Identity

// Context operations
func IdentityFromContext(ctx context.Context) *Identity
func AuthorizationFromContext(ctx context.Context) string

// Authorization
func Authorize(ctx context.Context, optionalId string,
    hook func(context.Context, *Identity) *Identity) (context.Context, *Identity, error)

// Checks
func IsAuthorized(ctx context.Context, optionalId string) bool

// Gin middleware helper
func PreAuthorize(c *gin.Context) *http.Request
```

**Usage Example:**

```go
// Create identity
identity := auth.NewIdentity("john@example.com", "John Doe", "john", "raw-data")

// Authorize with hook
ctx, identity, err := auth.Authorize(ctx, "admin", func(ctx context.Context, i *auth.Identity) *auth.Identity {
    if i.Email == "admin@example.com" {
        i.Authenticated("admin")
    }
    return i
})

// Check authorization
if auth.IsAuthorized(ctx, "admin") {
    // Grant access
}

// Gin middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        req := auth.PreAuthorize(c)
        ctx := req.Context()
        ctx, identity, err := auth.Authorize(ctx, "", nil)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        c.Set("identity", identity)
        c.Next()
    }
}
```

---

### 7. Search (`pkg/search/`)

Elasticsearch client with typed search abstractions.

**Key Types:**

```go
// Search results
type Hit[T interface{}, R interface{}] struct {
    Index  string  `json:"_index"`
    Id     string  `json:"_id"`
    Score  float64 `json:"_score"`
    Source T       `json:"_source"`
    Sort   *R      `json:"sort"`
}

type Hits[T interface{}, R interface{}] struct {
    Total Total       `json:"total"`
    Hits  []Hit[T, R] `json:"hits"`
}

type Records[T, R interface{}] struct {
    Hits Hits[T, R] `json:"hits"`
}

// Query types
type Query[R interface{}] struct {
    Size   int64              `json:"size,omitempty"`
    From   int64              `json:"from,omitempty"`
    Query  FunctionScoreQuery `json:"query,omitempty"`
    Source []string           `json:"_source,omitempty"`
    Sort   []any              `json:"sort,omitempty"`
    After  *R                 `json:"search_after,omitempty"`
}

// Document operations
type Doc map[types.ColumnName]interface{}

func NewDoc(id string) *Doc
func (d *Doc) SetId(id string) national.Sentence
func (d *Doc) Set(key types.ColumnName, value any, builtins maps.Set[types.ColumnName]) national.Sentence
func (d *Doc) SetNX(key types.ColumnName, value any)
func (d *Doc) Del(key types.ColumnName) (any, error)
```

**Usage Example:**

```go
// Create document
doc := search.NewDoc("user:123")
doc.SetNX(types.ColumnName("name"), "John Doe")
doc.SetNX(types.ColumnName("age"), 30)

// Build search query
query := search.Query[interface{}]{
    Size: 10,
    Query: search.FunctionScoreQuery{
        FunctionScore: search.FunctionScore{
            Query: search.BoolQuery{
                Bool: map[string]interface{}{
                    "must": []map[string]interface{}{
                        {"match": map[string]interface{}{"name": "John"}},
                    },
                },
            },
        },
    },
    Sort: []any{{"age": "desc"}},
}

// Execute search
results := es.Search(ctx, query, "users")
```

---

### 8. Common Types (`pkg/types/`)

Domain types and database utilities.

**Key Types:**

```go
// Database type
type DatabaseType string

const (
    MySQL      DatabaseType = "mysql"
    Postgres   DatabaseType = "postgres"
    Oracle     DatabaseType = "oracle"
    SQLite     DatabaseType = "sqlite"
    SQLServer  DatabaseType = "sqlserver"
    DB2        DatabaseType = "db2"
)

// Column name for type-safe references
type ColumnName string

// Parameter types
type ParameterType int

const (
    TypeBool ParameterType = iota
    TypeDouble
    TypeIdentity
    TypeInt
    TypeGroup
    TypeJson
    TypeString
    TypeText
    TypeArray
    TypeState
    TypeTimestamp
)

// Value types
type Float struct {
    Value *float64
    Valid bool
}

type IntCount struct {
    Value int64
    Valid bool
}

type ID struct {
    Value string
    Valid bool
}

// Match type for queries
type MatchType string

const (
    MatchTypeExact     MatchType = "exact"
    MatchTypeFuzzy     MatchType = "fuzzy"
    MatchTypePrefix    MatchType = "prefix"
    MatchTypeWildcard  MatchType = "wildcard"
)
```

**Key Functions:**

```go
// Column name quoting
func (c DatabaseType) Quote(i string) (o string)

// Value constructors
func NewFloat(f float64) Float
func NewIntCount(i int64) IntCount
func NewID(id string) ID

// Parameter validation
func (p ParameterType) Valid(value interface{}) error
func (p ParameterType) ColumnType() string
```

---

### 9. Time Utilities (`pkg/times/`)

Time formatting and timezone utilities.

**Key Functions:**

```go
// Format timestamps (milliseconds since epoch)
func FormatTs(ts int64) string
func FormatTsByLayout(ts int64, layout string) string

// UTC timestamp handling
func FormatTimestamp(timestamp int64, format string) string

// Timezone adjustments
func FixedTsByLocation(location *time.Location, time time.Time) types.IntCount
func FixedTs(time time.Time) types.IntCount

// Global offset
var UTCOffset int64 // UTC offset in microseconds
```

**Usage Example:**

```go
// Format to default layout "2006-01-02 15:04:05"
formatted := times.FormatTs(1731542400000) // "2024-11-14 00:00:00"

// Custom format
formatted := times.FormatTsByLayout(1731542400000, "2006/01/02") // "2024/11/02"

// UTC timestamp
utc := times.FormatTimestamp(1731542400000, "2006-01-02T15:04:05Z")
```

---

### 10. Language Utilities (`pkg/lang/`)

Go utility functions for common operations.

**Key Functions:**

```go
// Conditional operations
func If[T any, V any](expression bool, tf func(T) V, ff func(T) V, arg T) V
func Ternary[T any](expression bool, t T, f T) T
func Default[T any](v *T, d T) T
func ComputeIf[T interface{}](expression bool, tf func() T, ff func() T) T

// Pointer utilities
func Dup[T any](v T) *T
func OrDefault[T any](v *T, d *T) *T

// String utilities
func IsBlank(str string) bool
func JoinWith(sep string, parts ...*string) string
func TrimJoin(s string, parts ...*string) string
func Substring(p *string, start int, end int) string

// ID generation
func NewId() string
```

**Constants:**

```go
const Empty       = ""
const Space       = ' '
const SpaceString = " "
const Slash       = "/"
```

---

### 11. Map Utilities (`pkg/maps/`)

Enhanced map operations with type safety.

**Key Types:**

```go
type Map[K comparable, V interface{}] map[K]V

type Set[K comparable] interface {
    Has(k K) bool
}
```

**Key Functions:**

```go
// Map operations
func (m Map[K, V]) Has(k K) bool
func (m Map[K, V]) Get(k K) *V
func (m Map[K, V]) Put(k K, v V) Map[K, V]
func (m Map[K, V]) PutAll(a ...Map[K, V]) Map[K, V]
func (m Map[K, V]) Del(k K) Map[K, V]
func (m Map[K, V]) Merge(merge func(o, n V) (V, error), a ...Map[K, V]) (Map[K, V], error)
func (m Map[K, V]) PutIfAbsent(k K, f func(k K) V) V
func (m Map[K, V]) Equals(t Map[K, V], eq func(left, right V) bool) bool

// Utilities
func SetOf[K comparable](a ...K) Set[K]
func Keys[K comparable, V any](m Map[K, V]) lists.List[K]
func FromFn[K comparable, T, O any](arr []T, key func(T) K, val func(T) O) Map[K, O]
func Visitor[K comparable, V interface{}, P interface{}](
    m map[K]V, parent *P, tracer func(parent *P, key K) *P,
    drill func(value V) (map[K]V, bool), action func(parent *P, key K, value V))
```

---

### 12. List Utilities (`pkg/lists/`)

Enhanced list operations.

**Key Types:**

```go
type List[T any] []T
```

**Key Functions:**

```go
// Basic operations
func (l *List[T]) Length() int
func (l *List[T]) Empty() bool
func (l *List[T]) Add(a ...T) *List[T]
func (l *List[T]) Delete(a ...int) *List[T]

// Access
func (l *List[T]) Last() T
func (l *List[T]) FirstOrDefault(d *T) *T
func (l *List[T]) Sub(start, stop int) List[T]

// Functional
func (l *List[T]) For(fn func(t T))
func For[E any, V any](arr *[]E, m func(E) V) *[]V
func Collect[E any, O any](arr []E, f func(E) O) []O
func Filter[E any](arr []E, filter func(e E) bool) (List[E], List[E])

// String operations
func Join[E any](arr *[]E, separator string, f func(E) string) *string
func Of[T any](a ...T) List[T]
```

---

### 13. Error Handling (`pkg/errors/`)

Internationalized error definitions.

**Error Constants:**

```go
const (
    CompileExpressionError        i.Sentence = "Compile expression {{.express}} error"
    EvaluateExpressionError       i.Sentence = "Evaluate expression error: {{.error}}"
    ExpectedTypeButError          i.Sentence = "Type {{.expected}} is expected but got {{.actual}}"
    MissingTemplateError          i.Sentence = "Template is missing"
    MissingValueError             i.Sentence = "Missing required value"
    MissingAuthenticationToken    i.Sentence = "Missing authentication token"
    OverwriteInternalBuiltinError i.Sentence = "Internal builtin {{.builtin}} can't be overwritten"
    OverwriteBuiltinError         i.Sentence = "Builtin {{.builtin}} can't be overwritten"
    OverwriteIsForbiddenError     i.Sentence = "Overwrite {{.target}} of {{.type}} is not allowed"
    Unauthorized                  i.Sentence = "Unauthorized"
    UnexpectedTypeError           i.Sentence = "Got unexpected {{.type}}"
    UnexpectedValueError          i.Sentence = "Got unexpected {{.type}} {{.value}}"
    UnrecognizedError             i.Sentence = "Unrecognized {{.type}} {{.value}}"
    UnsupportedValueError         i.Sentence = "Unsupported {{.type}} {{.value}}"
    NotFoundError                 i.Sentence = "{{.type}} {{.value}} not found"
)
```

**Utility Functions:**

```go
func Yes(actors ...func()) *i.Message
func No(e error, argv ...any) *i.Message
```

---

### 14. Internationalization (`pkg/national/`)

i18n support with message templates.

**Key Types:**

```go
type Sentence string

type Localizer struct {
    // Internal implementation
}

type Message struct {
    // Internal implementation
}

var OkMessage Message
var ErrorMessage Message
```

**Key Functions:**

```go
// Message operations
func (s Sentence) LocalE(l *Localizer, logger *log.Logger, argv ...any) *Message
func (s Sentence) Local(l *Localizer, argv ...any) string

// Load translations
func LoadMessages(translations maps.Map[string, maps.Map[string, string]])
func LoadMessagesFromString(cfg string)

// Language detection
func WithLanguage(ctxt context.Context, req *http.Request, fallback string) context.Context
func Language(ctxt context.Context) string

// Get localizer
func Tr(ctxt context.Context) *Localizer
func FindOrCreateLocalizer(lang string) *Localizer
```

---

## Common Patterns

### Database Transaction with Cache Invalidation

```go
func UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
    return db.Tx(ctx, func(ctx context.Context, tx *gorm.DB) error {
        // Update in database
        if err := tx.Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
            return err
        }

        // Invalidate cache
        cache.Del(ctx, fmt.Sprintf("user:%s", id))

        return nil
    })
}
```

### Cache-Aside Pattern

```go
func GetUser(ctx context.Context, id string) (*User, error) {
    // Try cache
    cached, err := cache.Get(ctx, "user:"+id)
    if err == nil {
        return serde.Deserialize[User]([]byte(cached))
    }

    // Cache miss - query database
    user := &User{}
    if err := db.Query("users").Eq(types.ColumnName("id"), id).First(user).Error; err != nil {
        return nil, err
    }

    // Cache result
    if data, err := serde.Serialize(user); err == nil {
        cache.Set(ctx, "user:"+id, string(data), time.Hour)
    }

    return user, nil
}
```

### Authenticated Endpoint Pattern

```go
func GetUserProfile(c *gin.Context) {
    // Get identity from context (set by middleware)
    identity, _ := c.Get("identity")
    if identity == nil {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    id := identity.(*auth.Identity).UserName

    // Query user
    user, err := getUser(c.Request.Context(), id)
    if err != nil {
        c.JSON(404, gin.H{"error": "user not found"})
        return
    }

    c.JSON(200, user)
}
```

### Paginated Search Pattern

```go
func SearchUsers(ctx context.Context, query string, page, size int) (*SearchResult, error) {
    // Calculate offset
    offset := (page - 1) * size

    // Build query
    searchQuery := search.Query[interface{}]{
        Size: int64(size),
        From: int64(offset),
        Query: search.FunctionScoreQuery{
            FunctionScore: search.FunctionScore{
                Query: search.BoolQuery{
                    Bool: map[string]interface{}{
                        "should": []map[string]interface{}{
                            {"match": map[string]interface{}{"name": query}},
                            {"match": map[string]interface{}{"email": query}},
                        },
                    },
                },
            },
        },
        Sort: []any{{"created_at": "desc"}},
    }

    // Execute
    results := es.Search(ctx, searchQuery, "users")

    return &SearchResult{
        Total:   int(results.Hits.Total.Value),
        Page:    page,
        Size:    size,
        Results: results.Hits.Hits,
    }, nil
}
```

---

## Best Practices

1. **Use type-safe column references**: Always use `types.ColumnName("column_name")` instead of raw strings
2. **Leverage context propagation**: Pass context through all layers for tracing and i18n
3. **Use transactions**: Wrap multi-step database operations in `db.Tx()`
4. **Cache with expiration**: Always set TTL for cached values
5. **Handle errors internationally**: Use `national.Sentence().LocalE()` for user-facing errors
6. **Use generics**: Leverage generics in `serde` and `lists` for type safety
7. **Configure connection pools**: Adjust pool sizes based on workload
8. **Enable OpenTelemetry**: Setup tracing for production observability
9. **Validate inputs**: Use parameterized queries to prevent SQL injection
10. **Test with contexts**: Always pass context to database and cache operations

---

**Module:** `github.com/gantries/knife`
**License:** MIT
**Go Version:** 1.21+
