# About

`config` provides runtime configuration for environment variables.

# Usage

**cmd/basic_example.go**
```go
package main

import (
    "fmt"

    "github.com/drewscottt/config"
)

var _ config.Config = (*Config)(nil)

type Config struct {
    VarA int `config:"VAR_A"`
    VarB int `config:"VAR_B,default=20"`
}

func (c Config) Validate() error {
    return nil
}

func main() {
    c := Config{}
    err := config.Load(&c)
    if err != nil {
        panic(err)
    }
    fmt.Println("VarA", c.VarA)
    fmt.Println("VarB", c.VarB)
}
```

```bash
VAR_A=10 go run cmd/basic_example.go
```

```
VarA 10
VarB 20
```

# Suppported Types

Currently, native types `int`, `float64`, `string`, and `bool` are supported. Use `JSON` to support serialization of other data types (e.g. structs, slices, maps).

# Contribution

Just fork it, it will be easier.

