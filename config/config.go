package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	Addr     string `yaml:"addr"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	Debug    bool   `yaml:"debug"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

var defaultConfig = Config{
	Server: ServerConfig{
		Addr:     ":8080",
		CertFile: "",
		KeyFile:  "",
		Debug:    true,
	},
	Database: DatabaseConfig{
		Path: "data/data.db",
	},
	JWT: JWTConfig{
		Secret: "change-me-in-production",
	},
}

func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		data, _ := yaml.Marshal(&defaultConfig)
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = defaultConfig.JWT.Secret
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = defaultConfig.Database.Path
	}
	return &cfg, nil
}