package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaGenerateRes struct {
	Response string `json:"response"`
}

// OllamaEmbeddingRes 定义响应结构
type OllamaEmbeddingRes struct {
	Embedding []float32 `json:"embedding"`
}

func GetLocalEmbedding(text string) ([]float32, error) {
	url := "http://localhost:11434/api/embeddings"

	payload := map[string]string{
		"model":  "bge-m3", // 确保你本地运行了 ollama pull bge-m3
		"prompt": text,
	}

	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("请求 Ollama 失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res OllamaEmbeddingRes
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析 Ollama 响应失败: %w", err)
	}

	return res.Embedding, nil
}

func CallOllama(modelName string, prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"

	payload := map[string]any{
		"model":  modelName,
		"prompt": prompt,
		"stream": false, // 关闭流式输出，直接等待完整结果
	}

	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("请求 Ollama 失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res OllamaGenerateRes
	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("解析 Ollama 响应失败: %w", err)
	}

	return res.Response, nil
}
