package components

import (
	"github.com/charmbracelet/lipgloss"

	"r3f-trends/internal/adapter/driving/tui/styles"
)

type Header struct {
	width  int
	status string
}

func NewHeader() *Header {
	return &Header{
		status: "Ready",
	}
}

func (h *Header) SetWidth(width int) {
	h.width = width
}

func (h *Header) SetStatus(status string) {
	h.status = status
}

func (h *Header) View() string {
	title := styles.TitleStyle.Render("◈ R3F TREND COLLECTOR")
	statusText := styles.TrendMetaStyle.Render(h.status)
	spacer := lipgloss.NewStyle().Width(h.width - lipgloss.Width(title) - lipgloss.Width(statusText) - 4).Render(" ")

	return styles.HeaderStyle.Width(h.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, title, spacer, statusText),
	)
}
