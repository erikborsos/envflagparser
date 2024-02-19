package envflagparser_test

import (
	"os"
	"testing"
	"time"

	"github.com/erikborsos/envflagparser"
)

type Config struct {
	Host       string        `env:"HOST" flag:"host;localhost;Host address"`
	Timeout    time.Duration `env:"TIMEOUT" flag:"timeout;10s;Connection timeout"`
	EnableLogs bool          `env:"ENABLE_LOGS" flag:"enable-logs;false;Enable logging"`
}

func TestParseConfig(t *testing.T) {
	os.Setenv("HOST", "example.com")
	os.Setenv("TIMEOUT", "5s")
	os.Setenv("ENABLE_LOGS", "true")

	expectedConfig := Config{
		Host:       "example.com",
		Timeout:    5 * time.Second,
		EnableLogs: true,
	}

	var parsedConfig Config

	err := envflagparser.ParseConfig(&parsedConfig)
	if err != nil {
		t.Errorf("Error parsing config: %v", err)
	}

	if parsedConfig.Host != expectedConfig.Host {
		t.Errorf("Expected Host: %s, Got: %s", expectedConfig.Host, parsedConfig.Host)
	}
	if parsedConfig.Timeout != expectedConfig.Timeout {
		t.Errorf("Expected Timeout: %s, Got: %s", expectedConfig.Timeout, parsedConfig.Timeout)
	}
	if parsedConfig.EnableLogs != expectedConfig.EnableLogs {
		t.Errorf("Expected EnableLogs: %t, Got: %t", expectedConfig.EnableLogs, parsedConfig.EnableLogs)
	}
}
