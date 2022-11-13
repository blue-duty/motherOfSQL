package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB     DB
	Target Target
}

func ReadConfig(path string) (Config, error) {
	var config Config
	f, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
