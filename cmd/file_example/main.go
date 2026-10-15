package main

import (
	"fmt"
	"time"

	"github.com/drewscottt/config"
)

var _ config.Config = (*Config)(nil)

type Config struct {
	VarD *config.HotReload[int] `config:"VAR_D,default=10"`
}

func (c Config) Validate() error {
	return nil
}

func main() {
	c := Config{}
	err := config.Load(&c, config.AddConfigProvider(config.NewFileProvider(".env")))
	if err != nil {
		panic(err)
	}

	prevVarD, err := c.VarD.Get()
	if err != nil {
		panic(err)
	}
	fmt.Println("VarD", prevVarD)

	for {
		time.Sleep(5 * time.Second)
		newVarD, err := c.VarD.Get()
		if err != nil {
			panic(err)
		}
		if newVarD != prevVarD {
			fmt.Println("Updated VarD", newVarD)
		}
		prevVarD = newVarD
	}
}
