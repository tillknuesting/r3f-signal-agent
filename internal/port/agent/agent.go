package agent

import (
	"context"
)

type TopicSuggestion struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Score       float64  `json:"score"`
	Topics      []string `json:"topics"`
}

type Agent interface {
	Name() string
	Summarize(ctx context.Context, content string) (string, error)
	Suggest(ctx context.Context, trends []string, promptTemplate string) ([]TopicSuggestion, error)
}
