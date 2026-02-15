package markdown

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"r3f-trends/internal/domain/entity"
)

var ErrNotFound = errors.New("entity not found")

type TrendRepository struct {
	basePath string
}

func NewTrendRepository(basePath string) *TrendRepository {
	return &TrendRepository{basePath: basePath}
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

func (r *TrendRepository) Save(ctx context.Context, trend *entity.Trend) error {
	return r.SaveBatch(ctx, []*entity.Trend{trend})
}

func (r *TrendRepository) SaveBatch(ctx context.Context, trends []*entity.Trend) error {
	if len(trends) == 0 {
		return nil
	}

	profilePath := filepath.Join(r.basePath, "tech", "trends")
	if err := os.MkdirAll(profilePath, 0755); err != nil {
		return err
	}

	date := time.Now().Format("2006-01-02")
	filename := filepath.Join(profilePath, fmt.Sprintf("%s.md", date))

	existingTrends, err := r.loadFromFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	trendMap := make(map[string]*entity.Trend)
	for _, t := range existingTrends {
		trendMap[t.ID()] = t
	}
	for _, t := range trends {
		trendMap[t.ID()] = t
	}

	return r.saveToFile(filename, trendMap)
}

func (r *TrendRepository) FindByID(ctx context.Context, id string) (*entity.Trend, error) {
	trends, _, err := r.List(ctx, ListOptions{Limit: 10000})
	if err != nil {
		return nil, err
	}

	for _, t := range trends {
		if t.ID() == id {
			return t, nil
		}
	}

	return nil, ErrNotFound
}

func (r *TrendRepository) FindByDate(ctx context.Context, date string, opts ListOptions) ([]*entity.Trend, int, error) {
	profilePath := filepath.Join(r.basePath, "tech", "trends")
	filename := filepath.Join(profilePath, fmt.Sprintf("%s.md", date))

	trends, err := r.loadFromFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []*entity.Trend{}, 0, nil
		}
		return nil, 0, err
	}

	total := len(trends)

	start := opts.Offset
	if start > total {
		start = total
	}

	end := start + opts.Limit
	if end > total || opts.Limit == 0 {
		end = total
	}

	return trends[start:end], total, nil
}

func (r *TrendRepository) List(ctx context.Context, opts ListOptions) ([]*entity.Trend, int, error) {
	profilePath := filepath.Join(r.basePath, "tech", "trends")

	files, err := filepath.Glob(filepath.Join(profilePath, "*.md"))
	if err != nil {
		return nil, 0, err
	}

	var allTrends []*entity.Trend

	for _, file := range files {
		trends, err := r.loadFromFile(file)
		if err != nil {
			continue
		}
		allTrends = append(allTrends, trends...)
	}

	if opts.Source != "" {
		filtered := make([]*entity.Trend, 0)
		for _, t := range allTrends {
			if t.SourceID() == opts.Source {
				filtered = append(filtered, t)
			}
		}
		allTrends = filtered
	}

	total := len(allTrends)

	start := opts.Offset
	if start > total {
		start = total
	}

	end := start + opts.Limit
	if end > total || opts.Limit == 0 {
		end = total
	}

	return allTrends[start:end], total, nil
}

func (r *TrendRepository) Search(ctx context.Context, query string, opts SearchOptions) ([]*entity.Trend, int, error) {
	trends, _, err := r.List(ctx, ListOptions{Limit: 10000})
	if err != nil {
		return nil, 0, err
	}

	query = strings.ToLower(query)
	var results []*entity.Trend

	for _, t := range trends {
		if strings.Contains(strings.ToLower(t.Title()), query) ||
			strings.Contains(strings.ToLower(t.Summary()), query) {
			results = append(results, t)
		}
	}

	total := len(results)

	start := opts.Offset
	if start > total {
		start = total
	}

	end := start + opts.Limit
	if end > total || opts.Limit == 0 {
		end = total
	}

	return results[start:end], total, nil
}

func (r *TrendRepository) Update(ctx context.Context, trend *entity.Trend) error {
	return r.Save(ctx, trend)
}

func (r *TrendRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *TrendRepository) loadFromFile(filename string) ([]*entity.Trend, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content := string(data)

	frontmatterStart := strings.Index(content, "---")
	if frontmatterStart == -1 {
		return nil, fmt.Errorf("no frontmatter found")
	}

	secondFrontmatter := strings.Index(content[3:], "---")
	if secondFrontmatter == -1 {
		return nil, fmt.Errorf("no frontmatter end found")
	}

	afterFrontmatter := content[secondFrontmatter+6:]

	codeBlockStart := strings.Index(afterFrontmatter, "```json")
	if codeBlockStart == -1 {
		codeBlockStart = strings.Index(afterFrontmatter, "```")
	}

	var jsonContent string
	if codeBlockStart != -1 {
		afterCodeBlock := afterFrontmatter[codeBlockStart+3:]
		if strings.HasPrefix(strings.TrimSpace(afterCodeBlock), "json") {
			afterCodeBlock = strings.TrimSpace(afterCodeBlock)[4:]
		}
		codeBlockEnd := strings.Index(afterCodeBlock, "```")
		if codeBlockEnd == -1 {
			return []*entity.Trend{}, nil
		}
		jsonContent = strings.TrimSpace(afterCodeBlock[:codeBlockEnd])
	} else {
		jsonContent = strings.TrimSpace(afterFrontmatter)
	}

	if jsonContent == "" {
		return []*entity.Trend{}, nil
	}

	var trends []entity.TrendDTO
	if err := json.Unmarshal([]byte(jsonContent), &trends); err != nil {
		return nil, err
	}

	result := make([]*entity.Trend, len(trends))
	for i, dto := range trends {
		result[i] = entity.TrendFromDTO(&dto)
	}

	return result, nil
}

func (r *TrendRepository) saveToFile(filename string, trends map[string]*entity.Trend) error {
	var trendList []entity.TrendDTO
	for _, t := range trends {
		trendList = append(trendList, *t.ToDTO())
	}

	jsonData, err := json.MarshalIndent(trendList, "", "  ")
	if err != nil {
		return err
	}

	frontmatter := fmt.Sprintf(`---
date: %s
count: %d
---

`, time.Now().Format("2006-01-02"), len(trendList))

	var mdContent strings.Builder
	mdContent.WriteString(frontmatter)
	mdContent.WriteString("```json\n")
	mdContent.WriteString(string(jsonData))
	mdContent.WriteString("\n```\n")

	return os.WriteFile(filename, []byte(mdContent.String()), 0644)
}
