package main

import (
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func main() {
	spew.Dump(BigModelSearch("select * from d where a < 1 转成 dsl", "sk-2aab3fac28c64c5bb7919b8f94f9f518"))
}

func BigModelSearch(content string, appkey string) (resContent string, err error) {

	//调用大模型接口
	url := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

	// 构造请求体
	requestBody := map[string]interface{}{
		"model": "qwen-plus",
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "你是一个es程序员，能非常熟练的编写es查询语句。只需要回答es使用http方式查询的请求路径、请求方式、请求体即可，" +
					"最终以文本格式输出，不要使用Markdown语法。拒绝回答与es无关的问题。" +
					"参考格式为 ： 请求路径: /_search\n请求方式: GET\n请求体: \n{\n  \"query\": {\n    \"match_all\": {}\n  }\n}",
			},
			{
				"role":    "user",
				"content": content,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("JSON编码错误:", err)
		return
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+appkey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
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
		fmt.Println("JSON解析失败:", err)
		return
	}
	if len(responseBody.Choices) == 0 {
		err = errors.New("解析失败")
		return
	}

	resContent = responseBody.Choices[0].Message.Content

	return
}
