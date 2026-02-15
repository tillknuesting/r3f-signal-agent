package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"r3f-trends/internal/domain/entity"
	"r3f-trends/internal/domain/valueobject"
)

type HTTPCollector struct {
	client *http.Client
}

func New() *HTTPCollector {
	return &HTTPCollector{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *HTTPCollector) Type() valueobject.CollectorType {
	return valueobject.CollectorTypeHTTP
}

func (c *HTTPCollector) Validate(source *entity.Source) error {
	cfg := source.Config()
	if cfg["url"] == nil || cfg["url"] == "" {
		return fmt.Errorf("url is required")
	}
	return nil
}

func (c *HTTPCollector) Test(ctx context.Context, source *entity.Source) error {
	cfg := source.Config()
	url, _ := cfg["url"].(string)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

func (c *HTTPCollector) Collect(ctx context.Context, source *entity.Source) ([]*entity.Trend, error) {
	cfg := source.Config()

	url, _ := cfg["url"].(string)
	itemURL, _ := cfg["item_url"].(string)
	limit := 30
	if l, ok := cfg["limit"].(int); ok {
		limit = l
	}

	ids, err := c.fetchStoryIDs(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch story IDs: %w", err)
	}

	if len(ids) > limit {
		ids = ids[:limit]
	}

	trends := make([]*entity.Trend, 0, len(ids))
	fieldMapping := source.FieldMapping()

	for _, id := range ids {
		itemURLFilled := strings.ReplaceAll(itemURL, "{id}", strconv.Itoa(id))

		item, err := c.fetchItem(ctx, itemURLFilled)
		if err != nil {
			continue
		}

		trend := c.mapToTrend(item, source, fieldMapping)
		if trend != nil {
			trends = append(trends, trend)
		}
	}

	return trends, nil
}

func (c *HTTPCollector) fetchStoryIDs(ctx context.Context, url string) ([]int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ids []int
	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func (c *HTTPCollector) fetchItem(ctx context.Context, url string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item map[string]any
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}

	return item, nil
}

func (c *HTTPCollector) mapToTrend(item map[string]any, source *entity.Source, fieldMapping map[string]string) *entity.Trend {
	getField := func(apiField string) any {
		if mapped, ok := fieldMapping[apiField]; ok {
			apiField = mapped
		}
		return item[apiField]
	}

	id := ""
	switch v := getField("id").(type) {
	case float64:
		id = strconv.Itoa(int(v))
	case string:
		id = v
	}

	title, _ := getField("title").(string)
	url, _ := getField("url").(string)

	if title == "" {
		return nil
	}

	trend := entity.NewTrend(
		fmt.Sprintf("%s-%s", source.ID(), id),
		title,
		url,
	)
	trend.SetSource(source.Name())
	trend.SetSourceID(source.ID())

	if score, ok := getField("score").(float64); ok {
		trend.SetScore(int(score))
	}

	if author, ok := getField("author").(string); ok {
		trend.SetAuthor(author)
	}

	if ts, ok := getField("timestamp").(float64); ok {
		trend.SetTimestamp(time.Unix(int64(ts), 0))
	}

	for k, v := range item {
		if k != "id" && k != "title" && k != "url" && k != "score" && k != "by" && k != "time" {
			trend.SetMetadata(k, v)
		}
	}

	return trend
}
