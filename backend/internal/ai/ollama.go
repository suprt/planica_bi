package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient represents an Ollama API client
type OllamaClient struct {
	APIKey string
	APIURL string
	Model  string
	client *http.Client
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(apiKey, apiURL, model string) *OllamaClient {
	return &OllamaClient{
		APIKey: apiKey,
		APIURL: apiURL,
		Model:  model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateResponse represents a response from Ollama API
type GenerateResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

// ChatMessage represents a message in chat format
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
	Options  struct {
		Temperature float64 `json:"temperature"`
	} `json:"options,omitempty"`
	Temperature float64 `json:"temperature,omitempty"` // For OpenAI-compatible format
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Generate generates a response using Ollama API
func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("OLLAMA_API_KEY not set")
	}

	// Try multiple endpoint formats
	endpoints := []struct {
		url     string
		payload interface{}
		extract func(*http.Response) (string, error)
	}{
		// Local Ollama standard format (/generate)
		{
			url: fmt.Sprintf("%s/generate", c.APIURL),
			payload: map[string]interface{}{
				"model":  c.Model,
				"prompt": prompt,
				"stream": false,
				"options": map[string]interface{}{
					"temperature": 0.4,
				},
			},
			extract: func(resp *http.Response) (string, error) {
				var result GenerateResponse
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					return "", err
				}
				if result.Error != "" {
					return "", fmt.Errorf(result.Error)
				}
				return result.Response, nil
			},
		},
		// Local Ollama chat format (/chat)
		{
			url: fmt.Sprintf("%s/chat", c.APIURL),
			payload: ChatRequest{
				Model: c.Model,
				Messages: []ChatMessage{
					{
						Role:    "system",
						Content: "Ты анализируешь рекламные метрики и пишешь краткие выводы для клиента.",
					},
					{
						Role:    "user",
						Content: prompt,
					},
				},
				Stream: false,
			},
			extract: func(resp *http.Response) (string, error) {
				var result ChatResponse
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					return "", err
				}
				if result.Error.Message != "" {
					return "", fmt.Errorf(result.Error.Message)
				}
				if result.Message.Content != "" {
					return result.Message.Content, nil
				}
				return "", fmt.Errorf("empty response")
			},
		},
		// OpenAI-compatible format (/v1/chat/completions)
		{
			url: fmt.Sprintf("%s/v1/chat/completions", c.APIURL),
			payload: ChatRequest{
				Model: c.Model,
				Messages: []ChatMessage{
					{
						Role:    "system",
						Content: "Ты анализируешь рекламные метрики и пишешь краткие выводы для клиента.",
					},
					{
						Role:    "user",
						Content: prompt,
					},
				},
				Stream:      false,
				Temperature: 0.4,
			},
			extract: func(resp *http.Response) (string, error) {
				var result ChatResponse
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					return "", err
				}
				if result.Error.Message != "" {
					return "", fmt.Errorf(result.Error.Message)
				}
				if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
					return result.Choices[0].Message.Content, nil
				}
				if result.Message.Content != "" {
					return result.Message.Content, nil
				}
				return "", fmt.Errorf("empty response")
			},
		},
		// Try without /api prefix
		{
			url: func() string {
				baseURL := c.APIURL
				if baseURL[len(baseURL)-4:] == "/api" {
					baseURL = baseURL[:len(baseURL)-4]
				}
				return fmt.Sprintf("%s/v1/chat/completions", baseURL)
			}(),
			payload: ChatRequest{
				Model: c.Model,
				Messages: []ChatMessage{
					{
						Role:    "system",
						Content: "Ты анализируешь рекламные метрики и пишешь краткие выводы для клиента.",
					},
					{
						Role:    "user",
						Content: prompt,
					},
				},
				Stream:      false,
				Temperature: 0.4,
			},
			extract: func(resp *http.Response) (string, error) {
				var result ChatResponse
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					return "", err
				}
				if result.Error.Message != "" {
					return "", fmt.Errorf(result.Error.Message)
				}
				if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
					return result.Choices[0].Message.Content, nil
				}
				if result.Message.Content != "" {
					return result.Message.Content, nil
				}
				return "", fmt.Errorf("empty response")
			},
		},
	}

	var lastErr error
	for _, endpoint := range endpoints {
		// Prepare request
		payloadJSON, err := json.Marshal(endpoint.payload)
		if err != nil {
			lastErr = err
			continue
		}

		req, err := http.NewRequestWithContext(ctx, "POST", endpoint.url, bytes.NewBuffer(payloadJSON))
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

		// Make request
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
			continue
		}

		// Extract response
		result, err := endpoint.extract(resp)
		if err != nil {
			lastErr = err
			continue
		}

		return result, nil
	}

	return "", fmt.Errorf("failed to get response from Ollama API. Last error: %v", lastErr)
}
