package config

import (
	"io/ioutil"
	
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

func SaveConfig(cfg *Config) error {
	configPath := "config_dev/config.yml"
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.WithStack(err)
	}
	
	return ioutil.WriteFile(configPath, data, 0644)
}
