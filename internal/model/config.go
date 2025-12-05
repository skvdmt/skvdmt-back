package model

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	CONFIG_DIR           = "/etc"
	APP_DIR              = "skvdmt-back"
	POSTGRES_CONFIG_FILE = "postgres.yaml"
	SERVER_CONFIG_FILE   = "server.yaml"
)

// Config singleton app config
var Config *config

// PostgresConfig database connection config
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
}

// ServerConfig http server config
type ServerConfig struct {
	Port    uint16 `yaml:"port"`
	BaseUrl string `yaml:"base_url"`
}

// config wrapper configs
type config struct {
	Postgres *PostgresConfig
	Server   *ServerConfig
}

// LoadConfig load config files to singleton var Config
func LoadConfig() error {
	n := "model.LoadConfig"
	b, err := os.ReadFile(filepath.Join(CONFIG_DIR, APP_DIR, SERVER_CONFIG_FILE))
	if err != nil {
		return fmt.Errorf("%s %w", n, err)
	}
	s := &ServerConfig{}
	if err := yaml.Unmarshal(b, s); err != nil {
		return fmt.Errorf("%s %w", n, err)
	}
	b, err = os.ReadFile(filepath.Join(CONFIG_DIR, APP_DIR, POSTGRES_CONFIG_FILE))
	if err != nil {
		return fmt.Errorf("%s %w", n, err)
	}
	p := &PostgresConfig{}
	if err := yaml.Unmarshal(b, p); err != nil {
		return fmt.Errorf("%s %w", n, err)
	}
	Config = &config{
		Server:   s,
		Postgres: p,
	}
	return nil
}
