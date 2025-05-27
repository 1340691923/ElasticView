package config

import (
	"sync"
)

var (
	configInstance *Config
	configMutex    sync.RWMutex
)

func GetConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return configInstance
}

func SetConfigInstance(cfg *Config) {
	configMutex.Lock()
	defer configMutex.Unlock()
	configInstance = cfg
}
