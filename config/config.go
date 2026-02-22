package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         string `yaml:"port"`
	DatabasePath string `yaml:"database_path"`
}

const configPath = "./appdata/config.yaml"

func CheckIfExists() error {
	_, err := os.Stat(configPath)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	log.Printf("Config file %s doesn't exist, creating...", configPath)
	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}

	defaultConfig := Config{
		Port:         "8080",
		DatabasePath: "./appdata/database.db",
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
func LoadConfig() (Config, error) {
	var cfg Config

	file, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(file, &cfg)
	return cfg, err
}
