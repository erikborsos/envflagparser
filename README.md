# envflagparser

The `envflagparser` package provides functionality to parse configuration values from both environment variables and command-line flags into a provided struct. It offers the flexibility to prioritize environment variables over flag values.

## Features

- Parse configuration values from flags and environment variables into a provided struct.
- Prioritize environment variables over flag values, with an option to customize this behavior.
- Leverage reflection to dynamically set field values based on their types, making it convenient for configuring applications via flags or environment variables.

## Usage

1. Define a struct representing your configuration, with field tags specifying the corresponding environment variable names and flag names, along with default values and usage information.

```go
type Config struct {
    Host       string        `env:"HOST" flag:"host;localhost;Host address"`
    Port       string        `env:"PORT" flag:"port;localhost;Port"`
    Timeout    time.Duration `env:"TIMEOUT" flag:"timeout;20s;Connection timeout"`
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
    "time"

    "envflagparser"
)

type Config struct {
    Host       string        `env:"HOST" flag:"host;localhost;Host address"`
    Port       string        `env:"PORT" flag:"port;localhost;Port"`
    Timeout    time.Duration `env:"TIMEOUT" flag:"timeout;10s;Connection timeout"`
}

func main() {
    var config Config
    err := envflagparser.ParseConfig(&config)
    if err != nil {
        fmt.Printf("Error parsing config: %v\n", err)
        return
    }

    fmt.Printf("Host: %s\n", config.Host)
    fmt.Printf("Port: %s\n", config.Port)
    fmt.Printf("Timeout: %s\n", config.Timeout)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
