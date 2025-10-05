# Form Parser

A Go package for parsing HTTP request bodies into JSON, supporting various content types.

## Supported Content Types

- `application/json`: Parses JSON objects and arrays.
- `application/x-www-form-urlencoded`: Parses URL-encoded form data.
- `multipart/form-data`: Parses multipart form data (fields only, files are skipped).
- `text/plain`: Treats the body as plain text.

For unsupported content types, the raw body is returned.

## Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/oarkflow/form"
)

func main() {
    ctx := context.Background()
    contentType := "application/json"
    body := []byte(`{"name": "John", "age": 30}`)
    queryParams := map[string]string{"query": "test"}

    newCtx, jsonBytes, err := form.ParseBodyAsJSON(ctx, contentType, body, queryParams)
    if err != nil {
        panic(err)
    }

    // Access parsed data
    userCtx := form.UserContext(newCtx)
    fmt.Println(userCtx.Get("name")) // John
    fmt.Println(userCtx.Get("query")) // test

    // jsonBytes contains the JSON representation
    fmt.Println(string(jsonBytes)) // {"age":30,"name":"John","query":"test"}
}
```

## Context

The function enriches the context with a `Context` struct containing query parameters and parsed form data.

Use `form.UserContext(ctx)` to retrieve it.
