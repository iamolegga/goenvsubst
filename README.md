# goenvsubst

[![Go Reference](https://pkg.go.dev/badge/github.com/iamolegga/goenvsubst.svg)](https://pkg.go.dev/github.com/iamolegga/goenvsubst) [![Go Report Card](https://goreportcard.com/badge/github.com/iamolegga/goenvsubst)](https://goreportcard.com/report/github.com/iamolegga/goenvsubst) ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/iamolegga/goenvsubst/on-push-main.yml) [![Codacy Badge](https://app.codacy.com/project/badge/Coverage/62d6a131cbdb4c268069b279773deb2a)](https://app.codacy.com/gh/iamolegga/goenvsubst/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/62d6a131cbdb4c268069b279773deb2a)](https://app.codacy.com/gh/iamolegga/goenvsubst/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go package for recursively replacing environment variable references in Go data structures with their actual values from the environment.

## Features

- **Comprehensive Type Support**: Works with structs, slices, maps, arrays, and pointers
- **Nested Structures**: Handles deeply nested and complex data structures
- **Safe Pointer Handling**: Safely processes nil pointers without panics
- **In-Place Modification**: Modifies data structures in-place for efficiency
- **Environment Variable Format**: Uses `$VAR_NAME` format for variable references
- **Missing Variable Handling**: Replaces undefined or empty variables with empty strings
- **Zero Dependencies**: Pure Go implementation with no external dependencies

## Installation

```bash
go get github.com/iamolegga/goenvsubst
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"
    "github.com/iamolegga/goenvsubst"
)

func main() {
    // Set environment variable
    os.Setenv("DATABASE_URL", "postgres://localhost:5432/mydb")
    
    config := &struct {
        DatabaseURL string
        Debug       bool
    }{
        DatabaseURL: "$DATABASE_URL",
        Debug:       true,
    }
    
    err := goenvsubst.Do(config)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Database URL: %s\n", config.DatabaseURL)
    // Output: Database URL: postgres://localhost:5432/mydb
}
```

## Usage Examples

### Struct Fields

```go
type Config struct {
    Host     string
    Port     string
    Username string
    Password string
}

config := &Config{
    Host:     "$DB_HOST",
    Port:     "$DB_PORT", 
    Username: "$DB_USER",
    Password: "$DB_PASS",
}

goenvsubst.Do(config)
```

### Slices and Arrays

```go
// Slice
urls := &[]string{"$API_URL", "$BACKUP_URL", "https://static.example.com"}
goenvsubst.Do(urls)

// Array
servers := &[3]string{"$SERVER1", "$SERVER2", "$SERVER3"}
goenvsubst.Do(servers)
```

### Maps

```go
config := &map[string]string{
    "database_url": "$DATABASE_URL",
    "redis_url":    "$REDIS_URL",
    "api_key":      "$API_KEY",
}
goenvsubst.Do(config)
// Note: Map keys are never modified, only values
```

### Complex Nested Structures

```go
type DatabaseConfig struct {
    URL      string
    Username string
    Password string
}

type AppConfig struct {
    Database DatabaseConfig
    Services []string
    Env      map[string]string
}

config := &AppConfig{
    Database: DatabaseConfig{
        URL:      "$DATABASE_URL",
        Username: "$DB_USER",
        Password: "$DB_PASS",
    },
    Services: []string{"$SERVICE1", "$SERVICE2"},
    Env: map[string]string{
        "LOG_LEVEL": "$LOG_LEVEL",
        "DEBUG":     "$DEBUG_MODE",
    },
}

goenvsubst.Do(config)
```

### Pointers

```go
// Safe with nil pointers
var config *struct{ Value string }
goenvsubst.Do(&config) // No error, no operation

// Works with actual pointers
value := "$MY_VALUE"
ptr := &value
goenvsubst.Do(&ptr)
```

## Supported Data Types

| Type | Support | Notes |
|------|---------|-------|
| `string` | ✅ | Environment variables are substituted |
| `struct` | ✅ | All string fields are processed recursively |
| `slice` | ✅ | All elements are processed recursively |
| `array` | ✅ | All elements are processed recursively |
| `map` | ✅ | Only values are processed, keys remain unchanged |
| `pointer` | ✅ | Safely handles nil pointers |
| `int`, `bool`, etc. | ✅ | Non-string types are ignored (no substitution) |
| `interface{}` | ❌ | Not currently supported |

## Environment Variable Format

Environment variables should be referenced using the `$VAR_NAME` format:

```go
// Supported formats
"$DATABASE_URL"           // ✅ Full string replacement
"$MISSING_VAR"            // ✅ Undefined vars become empty strings

// Not supported formats
"${DATABASE_URL}"         // ❌ Braces not supported
"prefix-$API_KEY-suffix"  // ❌ Partial substitution not supported
```

## Error Handling

The `Do` function returns an error for future extensibility, but currently is designed to be robust:

```go
err := goenvsubst.Do(config)
if err != nil {
    log.Printf("Failed to substitute environment variables: %v", err)
    return
}
```

## Important Notes

- **In-Place Modification**: The function modifies the input data structure directly
- **Map Keys**: Only map values are processed, keys are never modified
- **Missing Variables**: Undefined or empty environment variables are replaced with empty strings
- **Thread Safety**: Safe for concurrent use (doesn't modify global state)
- **Nil Pointers**: Handled safely without causing panics
- **Type Safety**: Only string values are processed for substitution

## Testing

Run the tests:

```bash
go test -v
```

Run example tests:

```bash
go test -v -run Example
```

## Documentation

For detailed documentation and more examples, visit: [pkg.go.dev Documentation](https://pkg.go.dev/github.com/iamolegga/goenvsubst)

## License

MIT License - see [LICENSE](LICENSE) file for details.