package model

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// Название приложения.
	APP_NAME = "skvdmt-back"

	// Путь в директории конфигурации. (Добавляется директория с именем приложения).
	CONFIG_DIRECTORY_PROD = "/etc"
	CONFIG_DIRECTORY_DEV  = "./config"
	// Имя файла конфигурации.
	CONFIG_FILENAME_PROD = "prod.yaml"
	CONFIG_FILENAME_DEV  = "dev.yaml"
)

// Config Глобальная конфигурация.
var Config *MainConfig

// PostgresConfig Конфигурация соединения с postgres.
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
}

// ServerConfig Конфигурация HTTP сервера.
type ServerConfig struct {
	Port    uint16 `yaml:"port"`
	BaseUrl string `yaml:"base_url"`
}

// MainConfig Основная конфигурация.
type MainConfig struct {
	Postgres *PostgresConfig `yaml:"postgres"`
	Server   *ServerConfig   `yaml:"server"`
	Links    *LinksConfig    `yaml:"links"`
}

// LinksConfig Ссылки
type LinksConfig struct {
	Api           string `yaml:"api"`
	Documentation string `yaml:"documentation"`
}

// NewConfig Конфигурация.
func NewConfig() (*MainConfig, error) {
	Logs.Info.Info("configuration loading")
	f, err := os.ReadFile(filepath.Join(configDir(), configFilename()))
	if err != nil {
		return nil, err
	}
	c := &MainConfig{
		Postgres: &PostgresConfig{},
		Server:   &ServerConfig{},
		Links:    &LinksConfig{},
	}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, err
	}
	return c, nil
}

// configDir Директоия конфигурации.
func configDir() string {
	m, o := os.LookupEnv(MODE)
	if o && m == MODE_DEV {
		return CONFIG_DIRECTORY_DEV
	}
	return filepath.Join(CONFIG_DIRECTORY_PROD, APP_NAME)
}

// configFilename Имя файла конфигурации.
func configFilename() string {
	m, o := os.LookupEnv(MODE)
	if o && m == MODE_DEV {
		return CONFIG_FILENAME_DEV
	}
	return CONFIG_FILENAME_PROD
}
