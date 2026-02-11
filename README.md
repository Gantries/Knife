# knife

[![Go Report Card](https://goreportcard.com/badge/github.com/gantries/knife)](https://goreportcard.com/report/github.com/gantries/knife)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![CI](https://github.com/gantries/knife/actions/workflows/ci.yml/badge.svg)](https://github.com/gantries/knife/actions/workflows/ci.yml)

> A comprehensive Go utility library providing common tools and abstractions for building production-ready applications.

## Features

knife is a collection of well-tested, reusable Go packages organized into the following categories:

### Database

- **ORM Abstractions** - Unified database interface supporting MySQL, PostgreSQL, Oracle, SQLite, and SQL Server
- **Criteria Builder** - Fluent query building with type safety
- **Naming Strategies** - Configurable naming conventions for database entities

### Caching & Storage

- **Redis Integration** - High-level Redis client with HA support
- **Cache Abstractions** - Generic cache interfaces for flexibility

### Observability

- **OpenTelemetry Integration** - Tracing, metrics, and logging out of the box
- **Structured Logging** - Built on `log/slog` with OTel bridges

### Search

- **Elasticsearch Client** - Typed search abstractions for complex queries

### Utilities

- **Serialization** - Generic JSON/YAML (de)serialization with Go generics
- **Expression Evaluation** - Safe runtime expression evaluation
- **Internationalization** - i18n support with message templates
- **Maps & Lists** - Enhanced collection utilities
- **Time Utilities** - Timezone-aware formatting and conversion
- **Type Definitions** - Common domain types and converters

### Integration

- **Kubernetes Helpers** - Service account and secret access
- **Nacos Client** - Service discovery and configuration management
- **Authentication** - JWT and auth utilities

## Installation

```bash
go get github.com/gantries/knife@latest
```

## Quick Start

### Database ORM

```go
import "github.com/gantries/knife/pkg/orm"

// Configure database connection
db := orm.New(orm.Properties{
    Dialect:      types.Postgres,
    DSN:          "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable",
    TablePrefix:  "app_",
    SingularTable: true,
})

// Use with GORM
db.DB().AutoMigrate(&User{})
```

### Caching

```go
import "github.com/gantries/knife/pkg/cache"

// Create Redis cache
cache := cache.New(cache.Properties{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// Set and get values
cache.Set(ctx, "key", "value", time.Hour)
value, err := cache.Get(ctx, "key")
```

### Serialization

```go
import "github.com/gantries/knife/pkg/serde"

// Serialize/Deserialize with generics
data, err := serde.Serialize(user)
user, err := serde.Deserialize[User](data)

// YAML support
data, err := serde.SerializeYAML(config)
config, err := serde.DeserializeYAML[Config](data)
```

### Expression Evaluation

```go
import "github.com/gantries/knife/pkg/eval"

program, err := eval.Compile("user.Age > 18 && user.Active == true")
result, err := program.Eval(map[string]any{
    "user": map[string]any{"Age": 25, "Active": true},
})
// result: true
```

### OpenTelemetry

```go
import "github.com/gantries/knife/pkg/tel"

// Setup tracing
tp, err := tel.SetupTracer("my-service", "localhost:4317")
defer tel.ShutdownTracer(context.Background(), tp)

// Setup metrics
mp, err := tel.SetupMeter("my-service", "localhost:4317")
defer tel.ShutdownMeterProvider(context.Background(), mp)

// Setup logging
tel.SetupLoggersCreated(lp)
logger := tel.Logger("my-app")
```

### Time Utilities

```go
import "github.com/gantries/knife/pkg/times"

// Format timestamp
formatted := times.FormatTs(1731542400000)
// Output: "2024-11-14 00:00:00" (local timezone)

// Format with custom layout
formatted := times.FormatTsByLayout(1731542400000, "2006/01/02")
// Output: "2024/11/14"
```

## Documentation

For detailed documentation, see:

- [Godoc Reference](https://pkg.go.dev/github.com/gantries/knife) (once published)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)

## Project Status

knife is actively maintained and used in production environments. The library follows semantic versioning and aims for backward compatibility within major versions.

## Development

### Prerequisites

- Go 1.21 or later

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Building

```bash
# Build all packages
go build ./...
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

- Built with [GORM](https://gorm.io/)
- Uses [OpenTelemetry Go](https://opentelemetry.io/go/)
- Powered by [gin-gonic/gin](https://github.com/gin-gonic/gin)
- Search powered by [Elasticsearch](https://www.elastic.co/)

---

Made with ❤️ by the knife contributors
