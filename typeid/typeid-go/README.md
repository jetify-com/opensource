# TypeID Go
### A golang implementation of [TypeIDs](https://github.com/jetpack-io/typeid)
![License: Apache 2.0](https://img.shields.io/github/license/jetpack-io/typeid-go) [![Go Reference](https://pkg.go.dev/badge/go.jetpack.io/typeid.svg)](https://pkg.go.dev/go.jetpack.io/typeid)

TypeIDs are a modern, **type-safe**, globally unique identifier based on the upcoming
UUIDv7 standard. They provide a ton of nice properties that make them a great choice
as the primary identifiers for your data in a database, APIs, and distributed systems.
Read more about TypeIDs in their [spec](https://github.com/jetpack-io/typeid).

This particular implementation provides a go library for generating and parsing TypeIDs.

## Installation

To add this library as a dependency in your go module, run:

```bash
go get go.jetpack.io/typeid
```

## Usage
This library provides a go implementation of TypeID that optionally allows you
to define your own custom id types for added compile-time safety.

The statically typed version lives under the `typed` package. It makes it possible for
the go compiler itself to enforce type safety.

If you don't need compile-time safety, you can use the provided `typeid.TypeID` directly:
  
```go
import (
  "go.jetpack.io/typeid"
)

func example() {
  tid, _ := typeid.WithPrefix("user")
  fmt.Println(tid)
}
```

If you want compile-time safety, first define your own custom types with two steps:
1. Define your own struct and have it embed `typeid.TypeID`
2. Define an `AllowedPrefix` method.

You can now start using your custom-types as TypeIDs:

```go
import (
  "go.jetpack.io/typeid"
)

// To create a new id type, simply create a new struct, and have it embed TypeID:
type UserID struct {
	typeid.TypeID
}

// Then define AllowedPrefix(). In our case UserIDs use 'user' as a prefix
func (UserID) AllowedPrefix() string {
	return "user"
}

// That's it, you've now defined a subtype. Note that subtypes abide by the
// Subtype interface:
var _ typeid.Subtype = (*UserID)(nil)
```

And now use those types to generate TypeIDs:

```go
import (
  "go.jetpack.io/typeid/typed"
)

func example() {
  tid, _ := typeid.New[UserID]()
  fmt.Println(tid)
}
```

For the full documentation, see this package's [godoc](https://pkg.go.dev/go.jetpack.io/typeid).
