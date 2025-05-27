package config

import (
	"io/ioutil"
	
	"gopkg.in/yaml.v2"
)

func SaveConfig(cfg *Config) error {
	cfg.lock.Lock()
	defer cfg.lock.Unlock()
	
	configPath := "config_dev/config.yml"
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(configPath, data, 0644)
}
