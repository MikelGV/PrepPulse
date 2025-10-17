package tui

import "github.com/charmbracelet/lipgloss"

var (
    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#DB222A")).
        Padding(0, 1)
)

func (m Model) View() string {
    var output []string

    if m.err != nil {
        output = append(output, errorStyle.Render("Something went wrong opening file"))
    } else {
        output = append(output, m.table.Render())
    }
    return lipgloss.JoinVertical(lipgloss.Center, output...) 
}
