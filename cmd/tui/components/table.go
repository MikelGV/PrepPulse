package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rocketlaunchr/dataframe-go"
)

func CreateTable(df *dataframe.DataFrame, statsCol int, height, width int) string {
	/**
		I'm thinking of how to do it, here are a few thoughts that i have.
			1. 	the names(e.g. mm, etc, ...) are the buckets (or you create a bucket with its name).
			2. 	the cols data has to be matched with each bucket (so x data goes into a bucket).
			3. 	after all data is gathered in each bucket i have to display it
				in the x axis (this is the bins/buckets axis) and the y axis
				shows the frequency or count of data points that fall between each bin.
			4. 	all of this has to be shown in a popup table inside the dataframe file.
	**/
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("255"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("255")).
					Bold(true).
					Align(lipgloss.Center)
			case row%2 == 0:
				return lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("255")).
					Bold(true).
					Align(lipgloss.Center)
			default:
				return lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("255")).
					Bold(true).
					Align(lipgloss.Center)

			}
		})

	t.Row(
		fmt.Sprintf("%v", df),
	)

	content := lipgloss.JoinVertical(lipgloss.Center, t)

	return content
}
