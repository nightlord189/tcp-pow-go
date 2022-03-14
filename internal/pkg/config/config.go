package config

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"os"
)

// Config - Configuration for app (both client and server)
type Config struct {
	ServerHost            string `envconfig:"SERVER_HOST"`
	ServerPort            int    `envconfig:"SERVER_PORT"`
	HashcashZerosCount    int
	HashcashMaxIterations int
}

// Load - load config from file on path, after that get env-variables (override values from file)
func Load(path string) (*Config, error) {
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return &config, err
	}
	err = envconfig.Process("", &config)
	return &config, err
}
