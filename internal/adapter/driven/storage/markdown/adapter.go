package markdown

import (
	"context"

	"r3f-trends/internal/app/service"
	"r3f-trends/internal/domain/entity"
)

type TrendRepositoryAdapter struct {
	repo *TrendRepository
}

func NewTrendRepositoryAdapter(basePath string) *TrendRepositoryAdapter {
	return &TrendRepositoryAdapter{
		repo: NewTrendRepository(basePath),
	}
}

func (a *TrendRepositoryAdapter) SaveBatch(ctx context.Context, trends []*entity.Trend) error {
	return a.repo.SaveBatch(ctx, trends)
}

func (a *TrendRepositoryAdapter) List(ctx context.Context, opts service.ListOptions) ([]*entity.Trend, int, error) {
	return a.repo.List(ctx, ListOptions{
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Source: opts.Source,
		Date:   opts.Date,
	})
}

func (a *TrendRepositoryAdapter) FindByID(ctx context.Context, id string) (*entity.Trend, error) {
	return a.repo.FindByID(ctx, id)
}

func (a *TrendRepositoryAdapter) FindByDate(ctx context.Context, date string, opts service.ListOptions) ([]*entity.Trend, int, error) {
	return a.repo.FindByDate(ctx, date, ListOptions{
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Source: opts.Source,
		Date:   opts.Date,
	})
}

func (a *TrendRepositoryAdapter) Search(ctx context.Context, query string, opts service.SearchOptions) ([]*entity.Trend, int, error) {
	return a.repo.Search(ctx, query, SearchOptions{
		Limit:    opts.Limit,
		Offset:   opts.Offset,
		Sources:  opts.Sources,
		Tags:     opts.Tags,
		Starred:  opts.Starred,
		DateFrom: opts.DateFrom,
		DateTo:   opts.DateTo,
	})
}

func (a *TrendRepositoryAdapter) Update(ctx context.Context, trend *entity.Trend) error {
	return a.repo.Update(ctx, trend)
}

func (a *TrendRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return a.repo.Delete(ctx, id)
}
