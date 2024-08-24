package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SHWSAddress string  `yaml:"SHWAddress"`
	Addresses   Workers `yaml:"workerAddresses"`
	Logger      Logger  `yaml:"logger"`
	IsDebug     bool    `yaml:"isDebug"`
}

type Workers struct {
	Device string `yaml:"deviceWorkerAddress"`
	User   string `yaml:"userWorkerAddress"`
	Camera string `yaml:"cameraWorkerAddress"`
	Video  string `yaml:"videoStreamAddress"`
}

type Logger struct {
	Address string `yaml:"remote"`
	Path    string `yaml:"local"`
}

func New(path string) (*Config, error) {
	return initConfig(path)
}

func initConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, config)

	if err != nil {
		return nil, err
	}
	return &config, nil
}
