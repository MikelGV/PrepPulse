package components



/**
import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)
func CreateTable(df *dataframe.DataFrame) *table.Table {
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

    return t
}
**/
