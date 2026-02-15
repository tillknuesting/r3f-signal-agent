package glm5

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"r3f-trends/internal/app/service"
)

type Agent struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewAgent(apiKey, baseURL string) *Agent {
	if baseURL == "" {
		baseURL = "https://api.z.ai/v1"
	}
	return &Agent{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   "glm-5",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (a *Agent) Name() string {
	return "GLM-5"
}

func (a *Agent) Summarize(ctx context.Context, content string) (string, error) {
	prompt := fmt.Sprintf(`Summarize the following content in 2-3 sentences, highlighting key technical insights:

%s

Provide a concise summary:`, content)

	return a.chat(ctx, prompt)
}

func (a *Agent) Suggest(ctx context.Context, trends []string, promptTemplate string) ([]service.TopicSuggestion, error) {
	var prompt string
	if promptTemplate != "" {
		prompt = strings.ReplaceAll(promptTemplate, "{{topics}}", strings.Join(trends, "\n"))
	} else {
		prompt = fmt.Sprintf(`Based on these trending tech topics, suggest 3 blog post ideas that would fit a blog about cloud engineering, Go, Rust, and AI.

Topics:
%s

For each suggestion, provide a title and brief description in JSON format:
[
  {"title": "...", "description": "...", "score": 0.95}
]`, strings.Join(trends, "\n"))
	}

	response, err := a.chat(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return a.parseSuggestions(response)
}

func (a *Agent) chat(ctx context.Context, prompt string) (string, error) {
	req := chatRequest{
		Model: a.model,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (a *Agent) parseSuggestions(response string) ([]service.TopicSuggestion, error) {
	response = strings.TrimSpace(response)

	jsonStart := strings.Index(response, "[")
	jsonEnd := strings.LastIndex(response, "]")

	if jsonStart == -1 || jsonEnd == -1 {
		return a.parseTextSuggestions(response)
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var suggestions []service.TopicSuggestion
	if err := json.Unmarshal([]byte(jsonStr), &suggestions); err != nil {
		return a.parseTextSuggestions(response)
	}

	return suggestions, nil
}

func (a *Agent) parseTextSuggestions(response string) ([]service.TopicSuggestion, error) {
	suggestions := []service.TopicSuggestion{}

	lines := strings.Split(response, "\n")
	current := service.TopicSuggestion{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "##") {
			if current.Title != "" {
				suggestions = append(suggestions, current)
				current = service.TopicSuggestion{}
			}
			current.Title = strings.TrimLeft(line, "# ")
		} else if current.Title != "" && line != "" {
			current.Description += line + " "
		}
	}

	if current.Title != "" {
		suggestions = append(suggestions, current)
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, service.TopicSuggestion{
			Title:       "AI-Generated Suggestion",
			Description: response,
			Score:       0.8,
		})
	}

	return suggestions, nil
}
