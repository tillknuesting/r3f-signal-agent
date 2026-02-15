package storage

import (
	"context"

	"r3f-trends/internal/domain/entity"
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
	Save(ctx context.Context, trend *entity.Trend) error
	SaveBatch(ctx context.Context, trends []*entity.Trend) error
	FindByID(ctx context.Context, id string) (*entity.Trend, error)
	FindByDate(ctx context.Context, date string, opts ListOptions) ([]*entity.Trend, int, error)
	List(ctx context.Context, opts ListOptions) ([]*entity.Trend, int, error)
	Search(ctx context.Context, query string, opts SearchOptions) ([]*entity.Trend, int, error)
	Update(ctx context.Context, trend *entity.Trend) error
	Delete(ctx context.Context, id string) error
}

type SourceRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Source, error)
	FindAll(ctx context.Context) ([]*entity.Source, error)
	FindByProfile(ctx context.Context, profileName string) ([]*entity.Source, error)
	Save(ctx context.Context, source *entity.Source) error
	Delete(ctx context.Context, id string) error
}

type ProfileRepository interface {
	FindByName(ctx context.Context, name string) (*entity.Profile, error)
	FindAll(ctx context.Context) ([]*entity.Profile, error)
	FindActive(ctx context.Context) (*entity.Profile, error)
	Save(ctx context.Context, profile *entity.Profile) error
	Delete(ctx context.Context, name string) error
	SetActive(ctx context.Context, name string) error
}
