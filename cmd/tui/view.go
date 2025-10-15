package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
    var output []string
    output = append(output, m.table.Render())
    return lipgloss.JoinVertical(lipgloss.Center, output...) 
}
