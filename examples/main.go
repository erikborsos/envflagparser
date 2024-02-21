package main

import (
	"fmt"
	"os"
	"time"

	"github.com/erikborsos/envflagparser"
)

type MyConfig struct {
	Port    int           `env:"PORT" flag:"port" default:"8080" usage:"Server port"`
	AppName string        `env:"NAME" flag:"name" default:"MyApp" usage:"App name"`
	Debug   bool          `env:"DEBUG" flag:"debug" default:"false" usage:"Enable debug logs"`
	Timeout time.Duration `env:"TIMEOUT" flag:"timeout" default:"10s" usage:"Connection timeout"`
}

func main() {
	config := &MyConfig{}
	err := envflagparser.ParseConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing config: %v\n", err)
	}

	fmt.Println("Port:", config.Port)
	fmt.Println("AppName:", config.AppName)
	fmt.Println("Debug:", config.Debug)
	fmt.Println("Timeout:", config.Timeout)
}
