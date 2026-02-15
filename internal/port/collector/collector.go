package collector

import (
	"context"

	"r3f-trends/internal/domain/entity"
	"r3f-trends/internal/domain/valueobject"
)

type Collector interface {
	Type() valueobject.CollectorType
	Collect(ctx context.Context, source *entity.Source) ([]*entity.Trend, error)
	Validate(source *entity.Source) error
	Test(ctx context.Context, source *entity.Source) error
}

type CollectorFactory interface {
	Create(collectorType valueobject.CollectorType) (Collector, error)
}
