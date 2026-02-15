package service

import (
	"context"

	"r3f-trends/internal/domain/entity"
)

type CollectorService interface {
	Collect(ctx context.Context, profile string, sourceIDs []string) (*entity.CollectionResult, error)
	CollectAll(ctx context.Context, profile string) (*entity.CollectionResult, error)
	GetJobStatus(ctx context.Context, jobID string) (*entity.CollectionJob, error)
	ListJobs(ctx context.Context) ([]*entity.CollectionJob, error)
}

type TrendService interface {
	List(ctx context.Context, opts ListOptions) ([]*entity.Trend, int, error)
	Get(ctx context.Context, id string) (*entity.Trend, error)
	GetByDate(ctx context.Context, date string, opts ListOptions) ([]*entity.Trend, int, error)
	Search(ctx context.Context, query string, opts SearchOptions) ([]*entity.Trend, int, error)
	Star(ctx context.Context, id string) error
	Unstar(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type SourceService interface {
	Get(ctx context.Context, id string) (*entity.Source, error)
	List(ctx context.Context) ([]*entity.Source, error)
	ListByProfile(ctx context.Context, profile string) ([]*entity.Source, error)
	Create(ctx context.Context, source *entity.Source) error
	Update(ctx context.Context, source *entity.Source) error
	Delete(ctx context.Context, id string) error
	Test(ctx context.Context, id string) error
}

type ProfileService interface {
	Get(ctx context.Context, name string) (*entity.Profile, error)
	List(ctx context.Context) ([]*entity.Profile, error)
	GetActive(ctx context.Context) (*entity.Profile, error)
	Create(ctx context.Context, profile *entity.Profile) error
	Update(ctx context.Context, profile *entity.Profile) error
	Delete(ctx context.Context, name string) error
	Activate(ctx context.Context, name string) error
}

type AgentService interface {
	Summarize(ctx context.Context, trendID string) (string, error)
	SummarizeContent(ctx context.Context, content string) (string, error)
	SuggestTopics(ctx context.Context, profile string) ([]TopicSuggestion, error)
}

type ListOptions struct {
	Limit  int
	Offset int
	Source string
	Date   string
}

type SearchOptions struct {
	Limit    int
	Offset   int
	Sources  []string
	Tags     []string
	Starred  *bool
	DateFrom string
	DateTo   string
}

type TopicSuggestion struct {
	Title       string
	Description string
	Score       float64
	Topics      []string
}
