package components

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"r3f-trends/internal/adapter/driving/tui/styles"
)

type TrendItem struct {
	ID        string
	TitleText string
	URL       string
	Summary   string
	Score     int
	Source    string
	SourceID  string
	Author    string
	Timestamp time.Time
	Starred   bool
}

func (t TrendItem) FilterValue() string {
	return t.TitleText
}

type trendDelegate struct{}

func (d trendDelegate) Height() int  { return 2 }
func (d trendDelegate) Spacing() int { return 1 }

func (d trendDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d trendDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	t, ok := item.(TrendItem)
	if !ok {
		return
	}

	title := t.TitleText
	if t.Starred {
		title = styles.TrendStarredStyle.Render("★ ") + t.TitleText
	}

	meta := fmt.Sprintf("%s • %d pts", t.Source, t.Score)
	if t.Author != "" {
		meta += fmt.Sprintf(" • @%s", t.Author)
	}

	if index == m.Index() {
		titleStr := styles.TrendSelectedStyle.Render("▸ " + title)
		metaStr := styles.TrendMetaStyle.Render("  " + meta)
		fmt.Fprintln(w, titleStr)
		fmt.Fprint(w, metaStr)
	} else {
		titleStr := styles.TrendTitleStyle.Render("  " + title)
		metaStr := styles.TrendMetaStyle.Render("  " + meta)
		fmt.Fprintln(w, titleStr)
		fmt.Fprint(w, metaStr)
	}
}

type TrendsList struct {
	list    list.Model
	width   int
	height  int
	focused bool
}

func NewTrendsList() *TrendsList {
	l := list.New([]list.Item{}, trendDelegate{}, 0, 0)
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Title = "TRENDS"

	return &TrendsList{
		list:    l,
		focused: true,
	}
}

func (t *TrendsList) SetSize(width, height int) {
	t.width = width
	t.height = height
	t.list.SetSize(width, height)
}

func (t *TrendsList) SetFocused(focused bool) {
	t.focused = focused
}

func (t *TrendsList) Focused() bool {
	return t.focused
}

func (t *TrendsList) SetTrends(trends []TrendItem) {
	items := make([]list.Item, len(trends))
	for i, trend := range trends {
		items[i] = trend
	}
	t.list.SetItems(items)
}

func (t *TrendsList) SelectedTrend() *TrendItem {
	if item, ok := t.list.SelectedItem().(TrendItem); ok {
		return &item
	}
	return nil
}

func (t *TrendsList) CursorUp() {
	t.list.CursorUp()
}

func (t *TrendsList) CursorDown() {
	t.list.CursorDown()
}

func (t *TrendsList) View() string {
	return t.list.View()
}

type DetailView struct {
	width  int
	height int
	trend  *TrendItem
}

func NewDetailView() *DetailView {
	return &DetailView{}
}

func (d *DetailView) SetSize(width, height int) {
	d.width = width
	d.height = height
}

func (d *DetailView) SetTrend(trend *TrendItem) {
	d.trend = trend
}

func (d *DetailView) View() string {
	if d.trend == nil {
		return styles.BoxStyle.
			Width(d.width).
			Height(d.height).
			Render(styles.TrendMetaStyle.Render("Select a trend to view details"))
	}

	var b strings.Builder

	title := styles.TrendTitleStyle.Render(d.trend.TitleText)
	b.WriteString(title)
	b.WriteString("\n\n")

	meta := fmt.Sprintf("Source: %s\nScore: %d", d.trend.Source, d.trend.Score)
	if d.trend.Author != "" {
		meta += fmt.Sprintf("\nAuthor: @%s", d.trend.Author)
	}
	meta += fmt.Sprintf("\nTime: %s", d.trend.Timestamp.Format("2006-01-02 15:04"))
	b.WriteString(styles.TrendMetaStyle.Render(meta))

	if d.trend.Summary != "" {
		b.WriteString("\n\n")
		b.WriteString(d.trend.Summary)
	}

	b.WriteString("\n\n")
	url := styles.TrendSourceStyle.Render(d.trend.URL)
	b.WriteString(url)

	return styles.BoxStyle.
		Width(d.width).
		Height(d.height).
		Render(b.String())
}
