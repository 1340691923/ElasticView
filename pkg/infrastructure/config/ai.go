package config

import (
	"sync"
)

type AI struct {
	BigModeKey  string `json:"bigModeKey"`  // Qwen API key
	OpenAIKey   string `json:"openAIKey"`   // OpenAI API key
	DeepSeekKey string `json:"deepSeekKey"` // DeepSeek API key
}

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
