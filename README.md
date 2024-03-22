# envflagparser

The `envflagparser` package provides functionality to parse configuration values from both environment variables and command-line flags into a provided struct. It offers the flexibility to prioritize environment variables over flag values.

```sh
go get github.com/erikborsos/envflagparser
```
## Features

- Parse configuration values from flags and environment variables into a provided struct.
- Prioritize environment variables over flag values, with an option to customize this behavior.
- Leverage reflection to dynamically set field values based on their types, making it convenient for configuring applications via flags or environment variables.

## Usage

1. Define a struct representing your configuration, with field tags specifying the corresponding environment variable names and flag names, along with default values and usage information.

```go
type Config struct {
	Port    int           `env:"PORT" flag:"port" default:"8080" usage:"Server port"`
	AppName string        `env:"NAME" flag:"name" default:"MyApp" usage:"App name"`
	Debug   bool          `env:"DEBUG" flag:"debug" default:"false" usage:"Enable debug logs"`
	Timeout time.Duration `env:"TIMEOUT" flag:"timeout" default:"10s" usage:"Connection timeout"`
}
```

2. Call the `ParseConfig` function and pass a pointer to an instance of your configuration struct.

```go
var config &Config
err := envflagparser.ParseConfig(config)
if err != nil {
    // Handle error
}
```

3. Optionally, customize the behavior of the parser by modifying package-level variables such as `PrioritiseEnv` and `PrintErrorUsage`.

```go
envflagparser.PrioritiseEnv = false // Flags take precedence over environment variables
envflagparser.PrintErrorUsage = true // Include usage information in error messages
```

## Example

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/erikborsos/envflagparser"
)

type Config struct {
	Port    int           `env:"PORT" flag:"port" default:"8080" usage:"Server port"`
	AppName string        `env:"NAME" flag:"name" default:"MyApp" usage:"App name"`
	Debug   bool          `env:"DEBUG" flag:"debug" default:"false" usage:"Enable debug logs"`
	Timeout time.Duration `env:"TIMEOUT" flag:"timeout" default:"10s" usage:"Connection timeout"`
}

func main() {
	config := &Config{}
	err := envflagparser.ParseConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing config: %v\n", err)
	}

	fmt.Println("Port:", config.Port)
	fmt.Println("AppName:", config.AppName)
	fmt.Println("Debug:", config.Debug)
	fmt.Println("Timeout:", config.Timeout)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
