package chrome

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"

	"r3f-trends/internal/domain/entity"
	"r3f-trends/internal/domain/valueobject"
)

type ChromeCollector struct {
	mu       sync.Mutex
	timeout  time.Duration
	headless bool
}

func New() *ChromeCollector {
	return &ChromeCollector{
		timeout:  30 * time.Second,
		headless: true,
	}
}

func (c *ChromeCollector) Type() valueobject.CollectorType {
	return valueobject.CollectorTypeChrome
}

func (c *ChromeCollector) Validate(source *entity.Source) error {
	cfg := source.Config()
	if cfg["url"] == nil || cfg["url"] == "" {
		return fmt.Errorf("url is required")
	}
	return nil
}

func (c *ChromeCollector) Test(ctx context.Context, source *entity.Source) error {
	cfg := source.Config()
	url, _ := cfg["url"].(string)

	allocCtx, cancel := c.createContext(ctx)
	defer cancel()

	var result string
	return chromedp.Run(allocCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Title(&result),
	)
}

func (c *ChromeCollector) Collect(ctx context.Context, source *entity.Source) ([]*entity.Trend, error) {
	cfg := source.Config()

	url, _ := cfg["url"].(string)
	waitSelector, _ := cfg["wait_selector"].(string)
	containerSel, _ := cfg["container_selector"].(string)
	titleSel, _ := cfg["title_selector"].(string)
	linkSel, _ := cfg["link_selector"].(string)
	descSel, _ := cfg["description_selector"].(string)

	allocCtx, cancel := c.createContext(ctx)
	defer cancel()

	var results []map[string]string

	jsScript := c.buildExtractionScript(containerSel, titleSel, linkSel, descSel)

	err := chromedp.Run(allocCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		c.waitForElement(waitSelector),
		chromedp.Evaluate(jsScript, &results),
	)

	if err != nil {
		return nil, fmt.Errorf("chrome collection failed: %w", err)
	}

	trends := make([]*entity.Trend, 0, len(results))
	fieldMapping := source.FieldMapping()

	for i, item := range results {
		trend := c.mapToTrend(item, source, fieldMapping, i)
		if trend != nil {
			trends = append(trends, trend)
		}
	}

	return trends, nil
}

func (c *ChromeCollector) createContext(ctx context.Context) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", c.headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	timeoutCtx, timeoutCancel := context.WithTimeout(allocCtx, c.timeout)

	return timeoutCtx, func() {
		timeoutCancel()
		cancel()
	}
}

func (c *ChromeCollector) waitForElement(selector string) chromedp.Action {
	if selector == "" {
		return chromedp.Sleep(500 * time.Millisecond)
	}
	return chromedp.WaitVisible(selector)
}

func (c *ChromeCollector) buildExtractionScript(container, title, link, desc string) string {
	return fmt.Sprintf(`
	(function() {
		var results = [];
		var containers = document.querySelectorAll('%s');
		
		for (var i = 0; i < containers.length; i++) {
			var item = {};
			var container = containers[i];
			
			%s
			%s
			%s
			
			results.push(item);
		}
		
		return results;
	})()
	`,
		container,
		c.buildFieldScript("title", title),
		c.buildFieldScript("link", link, "href"),
		c.buildFieldScript("description", desc),
	)
}

func (c *ChromeCollector) buildFieldScript(name, selector string, attr ...string) string {
	if selector == "" {
		return ""
	}

	attribute := "innerText"
	if len(attr) > 0 {
		attribute = attr[0]
	}

	return fmt.Sprintf(`
		var %sEl = container.querySelector('%s');
		if (%sEl) {
			item.%s = %sEl.%s || %sEl.getAttribute('%s') || '';
			item.%s = item.%s.toString().trim();
		}
	`, name, selector, name, name, name, attribute, name, attribute, name, name)
}

func (c *ChromeCollector) mapToTrend(item map[string]string, source *entity.Source, fieldMapping map[string]string, index int) *entity.Trend {
	getField := func(key string) string {
		if mapped, ok := fieldMapping[key]; ok {
			return item[mapped]
		}
		return item[key]
	}

	title := getField("title")
	if title == "" {
		return nil
	}

	id := fmt.Sprintf("%s-%d-%d", source.ID(), time.Now().Unix(), index)
	url := getField("url")
	if url == "" {
		url = getField("link")
	}

	trend := entity.NewTrend(id, title, url)
	trend.SetSource(source.Name())
	trend.SetSourceID(source.ID())

	if summary := getField("summary"); summary != "" {
		trend.SetSummary(summary)
	}

	if description := getField("description"); description != "" && trend.Summary() == "" {
		trend.SetSummary(description)
	}

	for k, v := range item {
		if v != "" {
			trend.SetMetadata(k, v)
		}
	}

	return trend
}
