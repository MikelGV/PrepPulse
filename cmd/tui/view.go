package tui

import (
	"time"

	"github.com/MikelGV/PrepPulse/cmd/tui/components"
	"github.com/charmbracelet/lipgloss"
)

var (
    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#DB222A")).
        Padding(0, 1)

    warningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#A0A0A0")).
        Padding(0, 1)
)

func (m Model) View() string {
    var output []string

    if m.err != nil {
        output = append(output, errorStyle.Render("Error: " +m.err.Error()))
    } else if !m.ready {
        start := time.Now()

        if time.Since(start) > 100*time.Millisecond {
            output = append(output, warningStyle.Render("Scanning directory..."))
        } else {
            output = append(output, warningStyle.Render("Initializing..."))
        } 
    } else {
        switch m.state {
            case ListState:
                output = append(output, components.RenderList(m.files, m.selected, m.height, m.width))
            case FileState:
                if m.df != nil {
                    output = append(output, components.RenderDf(m.df, m.width, m.height, m.scrollCol, m.scrollRow))
                }
        }
    }

    return lipgloss.JoinVertical(lipgloss.Center, output...) 
}
