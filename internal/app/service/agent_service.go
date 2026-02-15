package service

import (
	"context"
	"strings"
)

type TopicSuggestion struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Score       float64  `json:"score"`
	Topics      []string `json:"topics"`
}

type LLMAgent interface {
	Name() string
	Summarize(ctx context.Context, content string) (string, error)
	Suggest(ctx context.Context, trends []string, promptTemplate string) ([]TopicSuggestion, error)
}

type AgentService struct {
	agent    LLMAgent
	trendSvc *TrendService
}

func NewAgentService(agent LLMAgent, trendSvc *TrendService) *AgentService {
	return &AgentService{
		agent:    agent,
		trendSvc: trendSvc,
	}
}

func (s *AgentService) Summarize(ctx context.Context, trendID string) (string, error) {
	trend, err := s.trendSvc.Get(ctx, trendID)
	if err != nil {
		return "", err
	}

	content := trend.Title()
	if trend.Summary() != "" {
		content = trend.Summary()
	}

	return s.agent.Summarize(ctx, content)
}

func (s *AgentService) SummarizeContent(ctx context.Context, content string) (string, error) {
	return s.agent.Summarize(ctx, content)
}

func (s *AgentService) SuggestTopics(ctx context.Context, profile string) ([]TopicSuggestion, error) {
	trends, _, err := s.trendSvc.List(ctx, ListOptions{Limit: 50})
	if err != nil {
		return nil, err
	}

	var trendTitles []string
	for _, t := range trends {
		trendTitles = append(trendTitles, t.Title())
	}

	prompt := s.buildPrompt(profile)

	return s.agent.Suggest(ctx, trendTitles, prompt)
}

func (s *AgentService) buildPrompt(profile string) string {
	switch profile {
	case "tech":
		return `Based on these trending tech topics, suggest 3 blog post ideas that would fit a blog about cloud engineering, Go, Rust, and AI.

Topics:
{{topics}}

For each suggestion, provide a title and brief description in JSON format:
[
  {"title": "...", "description": "...", "score": 0.95}
]`
	case "finance":
		return `Based on these trending finance topics, suggest 3 blog post ideas about investing and market analysis.

Topics:
{{topics}}

For each suggestion, provide a title and brief description in JSON format:
[
  {"title": "...", "description": "...", "score": 0.95}
]`
	default:
		return `Based on these trending topics, suggest 3 blog post ideas.

Topics:
{{topics}}

For each suggestion, provide a title and brief description in JSON format:
[
  {"title": "...", "description": "...", "score": 0.95}
]`
	}
}

func (s *AgentService) ExtractKeywords(ctx context.Context, content string) ([]string, error) {
	summary, err := s.agent.Summarize(ctx, content)
	if err != nil {
		return nil, err
	}

	words := strings.Fields(strings.ToLower(summary))
	keywords := make([]string, 0)

	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"may": true, "might": true, "must": true, "shall": true, "can": true,
		"to": true, "of": true, "in": true, "for": true, "on": true,
		"with": true, "at": true, "by": true, "from": true, "as": true,
		"into": true, "through": true, "during": true, "before": true, "after": true,
		"above": true, "below": true, "between": true, "under": true, "again": true,
		"further": true, "then": true, "once": true, "this": true, "that": true,
		"these": true, "those": true, "and": true, "but": true, "or": true,
		"if": true, "because": true, "while": true, "although": true, "though": true,
	}

	for _, word := range words {
		word = strings.Trim(word, ".,!?;:\"'()[]{}")
		if len(word) > 2 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	return keywords, nil
}
