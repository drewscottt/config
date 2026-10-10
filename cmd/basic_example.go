package main

import (
    "fmt"

    "github.com/drewscottt/config"
)

var _ config.Config = (*Config)(nil)

type Config struct {
    VarA int `config:"VAR_A"`
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
    fmt.Println(c.VarA)
}
