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
	configDirectoryProd = "/etc"
	configDirectoryDev  = "./config"
	// Имя файла конфигурации.
	configFileNameProd = "config.yaml"
	configFileNameDev  = "config-dev.yaml"
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
}

// LoadConfig Загрузка конфигурации в глобальную переменную Config.
func LoadConfig() error {
	Logs.Info.Info("configuration loading")
	configDirectory := configDirectoryProd
	configFileName := configFileNameProd
	mode, ok := os.LookupEnv(MODE)
	if ok && mode == Dev {
		configDirectory = configDirectoryDev
		configFileName = configFileNameDev
	}
	d, err := os.ReadFile(filepath.Join(configDirectory, APP_NAME, configFileName))
	if err != nil {
		return err
	}
	cfg := &MainConfig{}
	if err := yaml.Unmarshal(d, cfg); err != nil {
		return err
	}
	Config = cfg
	return nil
}
