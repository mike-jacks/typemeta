# typemeta

A lightweight Go package for managing type-level metadata with a powerful code generation tool.

## Features

- Type-safe metadata registration and retrieval
- Thread-safe operations with mutex protection
- Code generation support for automatic metadata registration
- Rich set of utility functions for metadata access
- Comprehensive test coverage and benchmarks

## Installation

```bash
go get github.com/mike-jacks/typemeta
```

## Quick Start

```go
package main

import "github.com/mike-jacks/typemeta"

type User struct {
    Name string
    Age  int
}

func main() {
    // Register metadata for a type
    typemeta.Register[User]("table", "users")
    typemeta.Register[User]("plural", "users")

    // Retrieve metadata
    if value, ok := typemeta.Meta[User]("table"); ok {
        fmt.Println(value) // Output: users
    }

    // Must-get with panic on missing key
    table := typemeta.Must[User]("table")

    // List all metadata entries
    entries := typemeta.List()
}
```

## Code Generation

The package includes a code generator tool that can automatically register metadata based on struct comments.

1. Install the generator:

```bash
go install github.com/mike-jacks/typemeta/cmd/typemeta-gen@latest
```

2. Add metadata comments to your structs:

```go
// +typemeta:table=users
// +typemeta:plural=users
type User struct {
    Name string
    Age  int
}
```

3. Run the generator:

```bash
typemeta-gen -root=./your/project
```

## API Reference

### Core Functions

- `Register[T any](key, value string)`: Register metadata for a type
- `Meta[T any](key string) (string, bool)`: Retrieve metadata with existence check
- `Must[T any](key string) string`: Retrieve metadata, panics if not found
- `MustWithLog[T any](key string) string`: Retrieve metadata with logging, panics if not found
- `List() []Entry`: List all registered metadata entries

### Types

```go
type Entry struct {
    TypeName string
    Key      string
    Value    string
}
```

## Benchmarks

The package includes comprehensive benchmarks for all operations. Run them with:

```bash
go test -bench=. ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
