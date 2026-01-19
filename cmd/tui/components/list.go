package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("b7a67b")).
			Padding(0, 1)
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#8ecae6"))
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0"))
)

func RenderList(files []string, selected int, width, height int) string {
	if len(files) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#de083a")).Render("No files found")
	}

	var items []string
	start := max(0, selected-10)
	end := min(len(files), start+height-3)

	for i := start; i < end; i++ {
		name := files[i]

		if len(name) > width-4 {
			name = name[:width-7] + "..."
		}

		if i == selected {
			items = append(items, selectedStyle.Render(" >"+name))
		} else {
			items = append(items, normalStyle.Render(" "+name))
		}
	}

	content := strings.Join(items, "\n")
	title := titleStyle.Render("Files [enter] [v]iew [e]dit [d]elete")

	return lipgloss.JoinVertical(lipgloss.Left, title, content)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
