package entity

import (
	"time"
)

type Trend struct {
	id          string
	title       string
	url         string
	summary     string
	score       int
	author      string
	source      string
	sourceID    string
	category    string
	tags        []string
	timestamp   time.Time
	collectedAt time.Time
	starred     bool
	metadata    map[string]any
}

func NewTrend(id, title, url string) *Trend {
	return &Trend{
		id:          id,
		title:       title,
		url:         url,
		collectedAt: time.Now(),
		metadata:    make(map[string]any),
		tags:        []string{},
	}
}

func (t *Trend) ID() string               { return t.id }
func (t *Trend) Title() string            { return t.title }
func (t *Trend) URL() string              { return t.url }
func (t *Trend) Summary() string          { return t.summary }
func (t *Trend) Score() int               { return t.score }
func (t *Trend) Author() string           { return t.author }
func (t *Trend) Source() string           { return t.source }
func (t *Trend) SourceID() string         { return t.sourceID }
func (t *Trend) Category() string         { return t.category }
func (t *Trend) Tags() []string           { return t.tags }
func (t *Trend) Timestamp() time.Time     { return t.timestamp }
func (t *Trend) CollectedAt() time.Time   { return t.collectedAt }
func (t *Trend) Starred() bool            { return t.starred }
func (t *Trend) Metadata() map[string]any { return t.metadata }

func (t *Trend) SetSummary(s string)             { t.summary = s }
func (t *Trend) SetScore(s int)                  { t.score = s }
func (t *Trend) SetAuthor(a string)              { t.author = a }
func (t *Trend) SetSource(s string)              { t.source = s }
func (t *Trend) SetSourceID(id string)           { t.sourceID = id }
func (t *Trend) SetCategory(c string)            { t.category = c }
func (t *Trend) SetTags(tags []string)           { t.tags = tags }
func (t *Trend) SetTimestamp(ts time.Time)       { t.timestamp = ts }
func (t *Trend) SetStarred(s bool)               { t.starred = s }
func (t *Trend) SetMetadata(key string, val any) { t.metadata[key] = val }
func (t *Trend) AddTag(tag string)               { t.tags = append(t.tags, tag) }

func (t *Trend) ToDTO() *TrendDTO {
	return &TrendDTO{
		ID:          t.id,
		Title:       t.title,
		URL:         t.url,
		Summary:     t.summary,
		Score:       t.score,
		Author:      t.author,
		Source:      t.source,
		SourceID:    t.sourceID,
		Category:    t.category,
		Tags:        t.tags,
		Timestamp:   t.timestamp,
		CollectedAt: t.collectedAt,
		Starred:     t.starred,
		Metadata:    t.metadata,
	}
}

type TrendDTO struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	URL         string         `json:"url"`
	Summary     string         `json:"summary,omitempty"`
	Score       int            `json:"score,omitempty"`
	Author      string         `json:"author,omitempty"`
	Source      string         `json:"source"`
	SourceID    string         `json:"source_id"`
	Category    string         `json:"category,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
	Timestamp   time.Time      `json:"timestamp"`
	CollectedAt time.Time      `json:"collected_at"`
	Starred     bool           `json:"starred"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

func TrendFromDTO(dto *TrendDTO) *Trend {
	t := NewTrend(dto.ID, dto.Title, dto.URL)
	t.summary = dto.Summary
	t.score = dto.Score
	t.author = dto.Author
	t.source = dto.Source
	t.sourceID = dto.SourceID
	t.category = dto.Category
	t.tags = dto.Tags
	t.timestamp = dto.Timestamp
	t.collectedAt = dto.CollectedAt
	t.starred = dto.Starred
	t.metadata = dto.Metadata
	if t.metadata == nil {
		t.metadata = make(map[string]any)
	}
	return t
}
