package service

import (
	"context"
	"fmt"

	"r3f-trends/internal/domain/entity"
	"r3f-trends/internal/domain/valueobject"
)

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

type TrendRepository interface {
	SaveBatch(ctx context.Context, trends []*entity.Trend) error
	List(ctx context.Context, opts ListOptions) ([]*entity.Trend, int, error)
	FindByID(ctx context.Context, id string) (*entity.Trend, error)
	FindByDate(ctx context.Context, date string, opts ListOptions) ([]*entity.Trend, int, error)
	Search(ctx context.Context, query string, opts SearchOptions) ([]*entity.Trend, int, error)
	Update(ctx context.Context, trend *entity.Trend) error
	Delete(ctx context.Context, id string) error
}

type Collector interface {
	Type() valueobject.CollectorType
	Collect(ctx context.Context, source *entity.Source) ([]*entity.Trend, error)
	Validate(source *entity.Source) error
}

type CollectorService struct {
	trendRepo  TrendRepository
	collectors map[string]Collector
}

func NewCollectorService(trendRepo TrendRepository, collectors map[string]interface{}) *CollectorService {
	c := make(map[string]Collector)
	for k, v := range collectors {
		if col, ok := v.(Collector); ok {
			c[k] = col
		}
	}
	return &CollectorService{
		trendRepo:  trendRepo,
		collectors: c,
	}
}

func (s *CollectorService) Collect(ctx context.Context, profile string, sourceIDs []string, sources []*entity.Source) (*entity.CollectionResult, error) {
	result := &entity.CollectionResult{
		Trends: []*entity.Trend{},
		Errors: []string{},
	}

	sourceMap := make(map[string]*entity.Source)
	for _, src := range sources {
		sourceMap[src.ID()] = src
	}

	for _, sourceID := range sourceIDs {
		source, exists := sourceMap[sourceID]
		if !exists {
			result.Errors = append(result.Errors, fmt.Sprintf("source not found: %s", sourceID))
			continue
		}

		if !source.Enabled() {
			continue
		}

		collector, exists := s.collectors[source.Type()]
		if !exists {
			result.Errors = append(result.Errors, fmt.Sprintf("collector not found for type: %s", source.Type()))
			continue
		}

		trends, err := collector.Collect(ctx, source)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to collect from %s: %v", sourceID, err))
			continue
		}

		result.Trends = append(result.Trends, trends...)
	}

	if len(result.Trends) > 0 {
		if err := s.trendRepo.SaveBatch(ctx, result.Trends); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to save trends: %v", err))
		}
	}

	return result, nil
}

type TrendService struct {
	repo TrendRepository
}

func NewTrendService(repo TrendRepository) *TrendService {
	return &TrendService{repo: repo}
}

func (s *TrendService) List(ctx context.Context, opts ListOptions) ([]*entity.Trend, int, error) {
	return s.repo.List(ctx, opts)
}

func (s *TrendService) Get(ctx context.Context, id string) (*entity.Trend, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TrendService) GetByDate(ctx context.Context, date string, opts ListOptions) ([]*entity.Trend, int, error) {
	return s.repo.FindByDate(ctx, date, opts)
}

func (s *TrendService) Search(ctx context.Context, query string, opts SearchOptions) ([]*entity.Trend, int, error) {
	return s.repo.Search(ctx, query, opts)
}

func (s *TrendService) Star(ctx context.Context, id string) error {
	trend, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	trend.SetStarred(true)
	return s.repo.Update(ctx, trend)
}

func (s *TrendService) Unstar(ctx context.Context, id string) error {
	trend, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	trend.SetStarred(false)
	return s.repo.Update(ctx, trend)
}

func (s *TrendService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
