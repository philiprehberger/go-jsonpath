# go-jsonpath

[![CI](https://github.com/philiprehberger/go-jsonpath/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-jsonpath/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-jsonpath.svg)](https://pkg.go.dev/github.com/philiprehberger/go-jsonpath) [![License](https://img.shields.io/github/license/philiprehberger/go-jsonpath)](LICENSE)

JSONPath query and extraction for Go. Type-safe with generics, zero dependencies

## Installation

```bash
go get github.com/philiprehberger/go-jsonpath
```

## Usage

### Extract Values

```go
import "github.com/philiprehberger/go-jsonpath"

data := []byte(`{"user": {"name": "Alice", "age": 30}}`)

name, err := jsonpath.Get[string](data, "$.user.name")
// name = "Alice"

age, err := jsonpath.Get[float64](data, "$.user.age")
// age = 30
```

### Wildcards

```go
data := []byte(`{"users": [{"name": "Alice"}, {"name": "Bob"}]}`)

names, err := jsonpath.GetAll[string](data, "$.users[*].name")
// names = ["Alice", "Bob"]
```

### Set Values

```go
data := []byte(`{"name": "Alice"}`)

updated, err := jsonpath.Set(data, "$.name", "Bob")
// updated = {"name": "Bob"}
```

### Exists

```go
import jsonpath "github.com/philiprehberger/go-jsonpath"

exists, _ := jsonpath.Exists(data, "$.address.city")
// true
```

### Delete

```go
import jsonpath "github.com/philiprehberger/go-jsonpath"

result, _ := jsonpath.Delete(data, "$.temporary")
```

### Supported Syntax

| Syntax | Description |
|--------|-------------|
| `$` | Root object |
| `$.field` | Object field |
| `$.a.b.c` | Nested fields |
| `$[0]` | Array index |
| `$.items[*]` | Wildcard (all elements) |

## API

| Function | Description |
|----------|-------------|
| `Get[T](data, path)` | Extract single value as type T |
| `GetRaw(data, path)` | Extract value without type conversion |
| `GetAll[T](data, path)` | Extract all wildcard matches as type T |
| `Set(data, path, value)` | Set value at path, return modified JSON |
| `Exists(data, path)` | Check whether a path exists |
| `Delete(data, path)` | Remove value at path, return modified JSON |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
