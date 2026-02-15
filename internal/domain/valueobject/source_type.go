package valueobject

type SourceType string

const (
	SourceTypeAggregator SourceType = "aggregator"
	SourceTypeCommunity  SourceType = "community"
	SourceTypeNews       SourceType = "news"
	SourceTypeResearch   SourceType = "research"
	SourceTypeCode       SourceType = "code"
	SourceTypeBlog       SourceType = "blog"
)
