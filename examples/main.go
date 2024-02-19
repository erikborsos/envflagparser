package main

import (
	"fmt"
	"os"
	"time"

	"github.com/erikborsos/envflagparser"
)

type MyConfig struct {
	Port    int           `env:"PORT" flag:"port;8080;Server port"`
	AppName string        `env:"NAME" flag:"name;MyApp;App name"`
	Debug   bool          `env:"DEBUG" flag:"debug;false;Enable debug logs"`
	Timeout time.Duration `env:"TIMEOUT" flag:"timeout;10s;Connection timeout"`
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
