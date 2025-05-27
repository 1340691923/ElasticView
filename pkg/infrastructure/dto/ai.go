package dto

type SearchBigModeReq struct {
	SysType int    `json:"sys_type"`
	Content string `json:"content"`
}

type AIConfigReq struct {
	QwenEnabled     bool   `json:"qwen_enabled"`
	BigModeKey      string `json:"big_mode_key"`
	OpenAIEnabled   bool   `json:"openai_enabled"`
	OpenAIKey       string `json:"open_ai_key"`
	DeepSeekEnabled bool   `json:"deepseek_enabled"`
	DeepSeekKey     string `json:"deep_seek_key"`
}

type AIConfigRes struct {
	QwenEnabled     bool   `json:"qwen_enabled"`
	BigModeKey      string `json:"big_mode_key"`
	OpenAIEnabled   bool   `json:"openai_enabled"`
	OpenAIKey       string `json:"open_ai_key"`
	DeepSeekEnabled bool   `json:"deepseek_enabled"`
	DeepSeekKey     string `json:"deep_seek_key"`
}
