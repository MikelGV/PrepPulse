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
		output = append(output, errorStyle.Render("Error: "+m.err.Error()))
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
				output = append(output, components.RenderDf(m.df, m.width, m.height, m.cursorRow, m.cursorCol, m.scrollRow, m.scrollCol, m.filteredRows, m.editMode, m.editBuffer))
			}
			if m.filterMode {
				prompt := lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")).Render("/")
				input := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render(m.filterInput + "_")
				bar := lipgloss.JoinHorizontal(lipgloss.Left, prompt, input)
				output = append(output, bar)
			}
			if m.statsMode {
				if stat, ok := m.chacheStats[m.statsColumns]; ok {
					popupWidth := int(float64(m.width) * 0.75)
					if popupWidth < 50 {
						popupWidth = 50
					}

					popupHeight := int(float64(m.height) * 0.80)
					if popupHeight < 20 {
						popupHeight = 20
					}

					if popupWidth >= m.width-2 {
						popupWidth = m.width - 4
					}
					if popupHeight >= m.height-2 {
						popupHeight = m.height - 4
					}

					componentsStat := components.ColumnsStat{
						Name:         stat.Name,
						Type:         stat.Type,
						Min:          stat.Min,
						Max:          stat.Max,
						Mean:         stat.Mean,
						Median:       stat.Median,
						Unique:       stat.Unique,
						Missing:      stat.Missing,
						Buckets:      stat.Buckets,
						BucketLabels: stat.BucketLabels,
					}
					popup := components.RenderHistogramPopup(componentsStat, popupWidth, popupHeight)
					overlay := lipgloss.Place(
						m.width, m.height,
						lipgloss.Center, lipgloss.Center,
						popup,
						lipgloss.WithWhitespaceChars(" "),
					)
					output = append(output, overlay)
				}
			}
		case HistogramInputState:
			if m.df != nil {
				output = append(output, components.RenderDf(m.df, m.width, m.height, m.cursorRow, m.cursorCol, m.scrollRow, m.scrollCol, m.filteredRows, m.editMode, m.editBuffer))
			}

			popupWidth := int(float64(m.width) * 0.60)
			if popupWidth < 50 {
				popupWidth = 50
			}
			if popupWidth >= m.width-2 {
				popupWidth = m.width - 4
			}

			popup := components.RenderColumnInputPopup(m.histInputBuffer, m.histSuggestions, m.histSelectedSugg, popupWidth, m.height)
			overlay := lipgloss.Place(
				m.width, m.height,
				lipgloss.Center, lipgloss.Center,
				popup,
				lipgloss.WithWhitespaceChars(" "),
			)
			output = append(output, overlay)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, output...)
}
