package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"r3f-trends/internal/adapter/driving/tui/styles"
)

type Footer struct {
	width    int
	lastRun  string
	count    int
	starred  int
	helpText string
}

func NewFooter() *Footer {
	return &Footer{
		lastRun:  "Never",
		count:    0,
		starred:  0,
		helpText: "[c] Collect  [s] Star  [↑↓] Navigate  [q] Quit",
	}
}

func (f *Footer) SetWidth(width int) {
	f.width = width
}

func (f *Footer) SetStats(lastRun string, count, starred int) {
	f.lastRun = lastRun
	f.count = count
	f.starred = starred
}

func (f *Footer) View() string {
	stats := fmt.Sprintf("Last run: %s • %d trends • %d starred", f.lastRun, f.count, f.starred)

	left := styles.TrendMetaStyle.Render(stats)
	right := styles.HelpStyle.Render(f.helpText)
	spacer := lipgloss.NewStyle().Width(f.width - lipgloss.Width(left) - lipgloss.Width(right) - 4).Render(" ")

	return styles.FooterStyle.Width(f.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, left, spacer, right),
	)
}

type Sidebar struct {
	width   int
	sources []SourceItem
	cursor  int
	focused bool
}

type SourceItem struct {
	ID      string
	Name    string
	Enabled bool
	Icon    string
}

func NewSidebar() *Sidebar {
	return &Sidebar{
		width:   22,
		sources: []SourceItem{},
		cursor:  0,
		focused: false,
	}
}

func (s *Sidebar) SetSources(sources []SourceItem) {
	s.sources = sources
}

func (s *Sidebar) SetFocused(focused bool) {
	s.focused = focused
}

func (s *Sidebar) Focused() bool {
	return s.focused
}

func (s *Sidebar) CursorUp() {
	if s.cursor > 0 {
		s.cursor--
	}
}

func (s *Sidebar) CursorDown() {
	if s.cursor < len(s.sources)-1 {
		s.cursor++
	}
}

func (s *Sidebar) ToggleSource() {
	if s.cursor < len(s.sources) {
		s.sources[s.cursor].Enabled = !s.sources[s.cursor].Enabled
	}
}

func (s *Sidebar) SelectedSource() *SourceItem {
	if s.cursor < len(s.sources) {
		return &s.sources[s.cursor]
	}
	return nil
}

func (s *Sidebar) View() string {
	var b strings.Builder

	title := styles.TitleStyle.Render("SOURCES")
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(styles.TrendMetaStyle.Render(strings.Repeat("─", 18)))
	b.WriteString("\n")

	for i, source := range s.sources {
		icon := source.Icon
		if icon == "" {
			icon = "•"
		}

		prefix := " "
		if s.focused && i == s.cursor {
			prefix = "▸"
		}

		check := "○"
		if source.Enabled {
			check = "◉"
		}

		line := fmt.Sprintf("%s %s %s", prefix, check, source.Name)

		if s.focused && i == s.cursor {
			b.WriteString(styles.TrendSelectedStyle.Render(line))
		} else if source.Enabled {
			b.WriteString(styles.SourceActiveStyle.Render(line))
		} else {
			b.WriteString(styles.SourceInactiveStyle.Render(line))
		}
		b.WriteString("\n")
	}

	return styles.SidebarStyle.Render(b.String())
}
