package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary    = lipgloss.Color("#1F2937")
	ColorAccent     = lipgloss.Color("#2563EB")
	ColorBackground = lipgloss.Color("#FAFAFA")
	ColorSuccess    = lipgloss.Color("#10B981")
	ColorWarning    = lipgloss.Color("#F59E0B")
	ColorError      = lipgloss.Color("#EF4444")
	ColorMuted      = lipgloss.Color("#6B7280")
	ColorBorder     = lipgloss.Color("#E5E7EB")

	HackerNewsColor = lipgloss.Color("#FF6600")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#1F2937")).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorAccent).
			Padding(0, 1).
			MarginBottom(1)

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1).
			Width(20)

	SourceStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SourceActiveStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess).
				PaddingLeft(2)

	SourceInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				PaddingLeft(2)

	TrendItemStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	TrendSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)

	TrendStarredStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F59E0B")).
				Bold(true)

	TrendTitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	TrendMetaStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	TrendSourceStyle = lipgloss.NewStyle().
				Foreground(HackerNewsColor)

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 1).
			MarginTop(1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			PaddingLeft(1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 2)

	StatusSuccessStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess).
				Bold(true)

	StatusErrorStyle = lipgloss.NewStyle().
				Foreground(ColorError).
				Bold(true)

	StatusPendingStyle = lipgloss.NewStyle().
				Foreground(ColorWarning).
				Bold(true)
)

func Width(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width)
}

func Height(height int) lipgloss.Style {
	return lipgloss.NewStyle().Height(height)
}
