package plugin

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type QQConfig struct {
	AppID  string `yaml:"appid"`
	Secret string `yaml:"secret"`
	Token  string `yaml:"token"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type Config struct {
	QQ     QQConfig     `yaml:"qq"`
	Server ServerConfig `yaml:"server"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}

	return cfg, nil
}
