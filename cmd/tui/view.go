package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
    var output []string
    return lipgloss.JoinVertical(lipgloss.Center, output...) 
}
