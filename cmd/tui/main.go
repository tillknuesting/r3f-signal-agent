package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"r3f-trends/internal/adapter/driving/tui/components"
	"r3f-trends/internal/adapter/driving/tui/styles"
)

type model struct {
	apiClient   *APIClient
	header      *components.Header
	footer      *components.Footer
	sidebar     *components.Sidebar
	trendsList  *components.TrendsList
	detailView  *components.DetailView
	width       int
	height      int
	focusedPane int
	loading     bool
	collecting  bool
	lastError   string
}

type trendsLoadedMsg struct {
	trends []components.TrendItem
	total  int
	err    error
}

type sourcesLoadedMsg struct {
	sources []components.SourceItem
	err     error
}

type collectCompleteMsg struct {
	result *CollectResponse
	err    error
}

type starCompleteMsg struct {
	trendID string
	err     error
}

func initialModel() model {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	return model{
		apiClient:   NewAPIClient(apiURL),
		header:      components.NewHeader(),
		footer:      components.NewFooter(),
		sidebar:     components.NewSidebar(),
		trendsList:  components.NewTrendsList(),
		detailView:  components.NewDetailView(),
		focusedPane: 0,
		loading:     true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		loadTrends(m.apiClient),
		loadSources(m.apiClient),
	)
}

func loadTrends(api *APIClient) tea.Cmd {
	return func() tea.Msg {
		resp, err := api.GetTrends()
		if err != nil {
			return trendsLoadedMsg{err: err}
		}

		trends := make([]components.TrendItem, len(resp.Trends))
		for i, t := range resp.Trends {
			trends[i] = components.TrendItem{
				ID:        t.ID,
				TitleText: t.Title,
				URL:       t.URL,
				Summary:   t.Summary,
				Score:     t.Score,
				Source:    t.Source,
				SourceID:  t.SourceID,
				Author:    t.Author,
				Timestamp: t.Timestamp,
				Starred:   t.Starred,
			}
		}

		return trendsLoadedMsg{trends: trends, total: resp.Total}
	}
}

func loadSources(api *APIClient) tea.Cmd {
	return func() tea.Msg {
		resp, err := api.GetSources()
		if err != nil {
			return sourcesLoadedMsg{err: err}
		}

		sources := make([]components.SourceItem, len(resp.Sources))
		for i, s := range resp.Sources {
			sources[i] = components.SourceItem{
				ID:      s.ID,
				Name:    s.Name,
				Enabled: s.Enabled,
				Icon:    s.Display.Icon,
			}
		}

		return sourcesLoadedMsg{sources: sources}
	}
}

func collectTrends(api *APIClient) tea.Cmd {
	return func() tea.Msg {
		result, err := api.Collect()
		return collectCompleteMsg{result: result, err: err}
	}
}

func starTrend(api *APIClient, trendID string) tea.Cmd {
	return func() tea.Msg {
		err := api.StarTrend(trendID)
		return starCompleteMsg{trendID: trendID, err: err}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.focusedPane = (m.focusedPane + 1) % 3
			m.sidebar.SetFocused(m.focusedPane == 0)
			m.trendsList.SetFocused(m.focusedPane == 1)
		case "shift+tab":
			m.focusedPane = (m.focusedPane - 1 + 3) % 3
			m.sidebar.SetFocused(m.focusedPane == 0)
			m.trendsList.SetFocused(m.focusedPane == 1)
		case "up", "k":
			if m.focusedPane == 0 {
				m.sidebar.CursorUp()
			} else if m.focusedPane == 1 {
				m.trendsList.CursorUp()
				if t := m.trendsList.SelectedTrend(); t != nil {
					m.detailView.SetTrend(t)
				}
			}
		case "down", "j":
			if m.focusedPane == 0 {
				m.sidebar.CursorDown()
			} else if m.focusedPane == 1 {
				m.trendsList.CursorDown()
				if t := m.trendsList.SelectedTrend(); t != nil {
					m.detailView.SetTrend(t)
				}
			}
		case " ":
			if m.focusedPane == 0 {
				m.sidebar.ToggleSource()
			}
		case "c":
			if !m.collecting {
				m.collecting = true
				m.header.SetStatus("Collecting...")
				cmds = append(cmds, collectTrends(m.apiClient))
			}
		case "s":
			if m.focusedPane == 1 {
				if t := m.trendsList.SelectedTrend(); t != nil {
					cmds = append(cmds, starTrend(m.apiClient, t.ID))
				}
			}
		case "enter":
			if t := m.trendsList.SelectedTrend(); t != nil {
				m.detailView.SetTrend(t)
			}
		case "r":
			m.loading = true
			m.header.SetStatus("Refreshing...")
			cmds = append(cmds, loadTrends(m.apiClient))
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.header.SetWidth(m.width)
		m.footer.SetWidth(m.width)

		sidebarWidth := 24
		mainWidth := m.width - sidebarWidth - 2
		listWidth := mainWidth / 2
		detailWidth := mainWidth - listWidth - 2

		contentHeight := m.height - 6

		m.trendsList.SetSize(listWidth, contentHeight)
		m.detailView.SetSize(detailWidth, contentHeight)

	case trendsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.lastError = msg.err.Error()
			m.header.SetStatus("Error loading trends")
		} else {
			m.trendsList.SetTrends(msg.trends)
			m.footer.SetStats("Recently", msg.total, 0)
			m.header.SetStatus(fmt.Sprintf("Loaded %d trends", msg.total))
			m.lastError = ""
			if len(msg.trends) > 0 {
				m.detailView.SetTrend(&msg.trends[0])
			}
		}

	case sourcesLoadedMsg:
		if msg.err == nil {
			m.sidebar.SetSources(msg.sources)
		}

	case collectCompleteMsg:
		m.collecting = false
		if msg.err != nil {
			m.header.SetStatus("Collection failed")
			m.lastError = msg.err.Error()
		} else {
			m.header.SetStatus(fmt.Sprintf("Collected %d trends", msg.result.ItemsCount))
			cmds = append(cmds, loadTrends(m.apiClient))
		}

	case starCompleteMsg:
		if msg.err == nil {
			cmds = append(cmds, loadTrends(m.apiClient))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	header := m.header.View()
	footer := m.footer.View()

	sidebar := m.sidebar.View()
	trendsList := m.trendsList.View()
	detailView := m.detailView.View()

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		styles.Width(2).Render(""),
		trendsList,
		styles.Width(2).Render(""),
		detailView,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		mainContent,
		footer,
	)
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
