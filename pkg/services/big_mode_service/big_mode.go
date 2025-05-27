package big_mode_service

import (
	"bytes"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const DefaultSysContent = `你是一个高级程序员，能非常熟练的编写es,mysql,redis,clickhouse,postgres,mongo查询语句。" +
					"最终结果以文本格式输出，不要使用Markdown语法。拒绝回答与代码无关的问题。`

const DefaultSysMd = `你是一个高级程序员，能非常熟练的编写es,mysql,redis,clickhouse,postgres,mongo查询语句。" +
					"最终结果以Markdown语法格式输出，不要使用文本。拒绝回答与代码无关的问题。`

const CommonSysMd = ` 最终结果以Markdown语法格式输出，不要使用文本。。`

const CompatibleModeUrl = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

type BigMode struct {
	log         *logger.AppLogger
	qwenKey     string
	openaiKey   string
	deepseekKey string
}

func NewBigMode(log *logger.AppLogger, cfg *config.Config) *BigMode {
	return &BigMode{
		log:         log,
		qwenKey:     cfg.Ai.BigModeKey,
		openaiKey:   cfg.Ai.OpenAIKey,
		deepseekKey: cfg.Ai.DeepSeekKey,
	}
}

func (this *BigMode) BigModelSearch(sysContent, content string) (resContent string, err error) {

	if this.qwenKey == "" {
		err = errors.New("请修改config.yml中的bigModeKey,获取地址: https://bailian.console.aliyun.com/?apiKey=1#/api-key ")
		return
	}

	//调用大模型接口
	url := CompatibleModeUrl

	// 构造请求体
	requestBody := map[string]interface{}{
		"model": "qwen-plus",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": sysContent,
			},
			{
				"role":    "user",
				"content": content,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+this.qwenKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}

	type ResBody struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string      `json:"finish_reason"`
			Index        int         `json:"index"`
			Logprobs     interface{} `json:"logprobs"`
		} `json:"choices"`
		Object string `json:"object"`
		Usage  struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Created           int         `json:"created"`
		SystemFingerprint interface{} `json:"system_fingerprint"`
		Model             string      `json:"model"`
		Id                string      `json:"id"`
	}

	fmt.Println("响应状态:", resp.Status)
	var responseBody ResBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	if len(responseBody.Choices) == 0 {
		err = errors.New("解析失败")
		this.log.Error("err", zap.Error(err))
		return
	}

	resContent = responseBody.Choices[0].Message.Content

	return
}

func (this *BigMode) ChatGPTSearch(sysContent, content string) (resContent string, err error) {
	if this.openaiKey == "" {
		err = errors.New("请配置OpenAI API Key")
		return
	}
	
	url := "https://api.openai.com/v1/chat/completions"
	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": sysContent},
			{"role": "user", "content": content},
		},
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	
	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	
	req.Header.Set("Authorization", "Bearer "+this.openaiKey)
	req.Header.Set("Content-Type", "application/json")
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	
	type ResBody struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string      `json:"finish_reason"`
			Index        int         `json:"index"`
			Logprobs     interface{} `json:"logprobs"`
		} `json:"choices"`
		Object string `json:"object"`
		Usage  struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	
	var responseBody ResBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	if len(responseBody.Choices) == 0 {
		err = errors.New("解析失败")
		this.log.Error("err", zap.Error(err))
		return
	}
	
	resContent = responseBody.Choices[0].Message.Content
	
	return
}

func (this *BigMode) DeepSeekSearch(sysContent, content string) (resContent string, err error) {
	if this.deepseekKey == "" {
		err = errors.New("请配置DeepSeek API Key")
		return
	}
	
	url := "https://api.deepseek.com/v1/chat/completions"
	requestBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "system", "content": sysContent},
			{"role": "user", "content": content},
		},
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	
	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	
	req.Header.Set("Authorization", "Bearer "+this.deepseekKey)
	req.Header.Set("Content-Type", "application/json")
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	
	type ResBody struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string      `json:"finish_reason"`
			Index        int         `json:"index"`
			Logprobs     interface{} `json:"logprobs"`
		} `json:"choices"`
		Object string `json:"object"`
		Usage  struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	
	var responseBody ResBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		this.log.Error("err", zap.Error(err))
		err = errors.WithStack(err)
		return
	}
	if len(responseBody.Choices) == 0 {
		err = errors.New("解析失败")
		this.log.Error("err", zap.Error(err))
		return
	}
	
	resContent = responseBody.Choices[0].Message.Content
	
	return
}
